// Command gen generates go-swagger yaml and handlers from route data.
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
)

var (
	//go:embed data/swagger.goyaml
	ymlTmpl string

	//go:embed data/code.gogotmpl
	goTmpl string
)

func main() {
	cmdi.RegisterRoutes(nil)
	// TODO other routes
	rmap := common.PopulateRouteMap(common.AllRoutes)
	var rlist = make(common.Routes, 0, len(rmap))
	for _, r := range rmap {
		if len(r.Sinks) == 0 || len(r.Sinks[0].Name) == 0 {
			// skip
			continue
		}
		switch strings.Trim(r.Base, "/ ") {
		case "cmdInjection":
		default:
			// skip others for now
			continue
		}
		if len(r.Sinks) > 0 && len(r.Sinks[0].Name) > 0 {
			rlist = append(rlist, r)
		}
	}
	//sort so the generated code is stable
	sort.Sort(rlist)

	cmdDir, err := findSwagCmd()
	if err != nil {
		log.Fatal(err)
	}
	genYml, err := os.Create(filepath.Join(cmdDir, "swagger.yml"))
	if err != nil {
		log.Fatal(err)
	}
	defer genYml.Close()

	tdata := tmplData{
		GenNotice: "GENERATED CODE - DO NOT EDIT",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		GenCmd:    "go run ./cmd/go-swagger/regen/regen.go",
		Routes:    rlist,
	}
	tfuncs := template.FuncMap{
		"capital":         strings.Title,
		"assignStmt":      assignStmt,
		"routePkg":        routePkg,
		"routeIdentifier": routeIdentifier,
	}

	if err = generateYaml(tdata, tfuncs, genYml); err != nil {
		log.Fatal(err)
	}
	if err = runSwagger(); err != nil {
		log.Fatal(err)
	}
	swagPkg, err := findSwagPkg()
	if err != nil {
		log.Fatal(err)
	}

	genGo := filepath.Join(swagPkg, "generatedInit.go")
	g, err := os.Create(genGo)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	if err = generateCode(tdata, tfuncs, g); err != nil {
		log.Fatal(err)
	}
}

type tmplData struct {
	GenNotice string // kept separate from template so that automated stuff doesn't identify the template itself as generated.
	Timestamp string
	GenCmd    string
	Routes    common.Routes
}

func generateYaml(tdata tmplData, tfuncs template.FuncMap, w io.Writer) error {
	tmpl, err := template.New("yaml").Funcs(tfuncs).Parse(ymlTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, tdata)
}

func runSwagger() error {
	if false {
		exec.Command("go-swagger" /*... TODO */)
		panic("unimplemented")
	}
	return nil
}

func generateCode(td tmplData, tfuncs template.FuncMap, w io.Writer) error {
	tmpl, err := template.New("go").Funcs(tfuncs).Parse(goTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, td)
}

func assignStmt(r *common.Route, s *common.Sink, i string) string {
	m := sinkFn(i, s)
	return fmt.Sprintf("api.%s%sHandler = %s.%sHandlerFunc(%s)", routeIdentifier(r), m, routePkg(r), m, m)
}

//sink name, as used by swagger in url
func sinkName(s *common.Sink) string {
	fname := s.URL
	idx := strings.LastIndexByte(fname, '.')
	if idx < 0 {
		return exportIdentifier(fname)
	}
	return exportIdentifier(fname[idx+1:])
}

// name of wrapper around vulnerable function
func sinkFn(in string, s *common.Sink) string {
	return "Get" + strings.Title(in) + sinkName(s)
}

// generates package name swagger uses for route
// CmdInjection -> cmd_injection
func routePkg(r *common.Route) string {
	pkg := strings.ToLower(r.Base)
	// ignore 0th letter - start at 1
	j := 1
	for i := 1; i < len(r.Base); i++ {
		// fmt.Fprintf(os.Stderr, "i=%d j=%d n=%d N=%d\n", i, j, len(name), len(r.Name))
		if pkg[j] != r.Base[i] {
			//case changed, insert underscore (and advance 1)
			pkg = pkg[:j] + "_" + pkg[j:]
			j++
		}
		j++
	}
	pkg = strings.TrimLeft(pkg, " -/")
	return strings.TrimRight(pkg, " -/")
}

//return an identifier for the route, suitable for use in an exported function name
func routeIdentifier(r *common.Route) string {
	return exportIdentifier(r.Base)
}
func exportIdentifier(id string) string {
	id = strings.TrimLeft(id, "-./ _")
	id = strings.TrimRight(id, "-./ +")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ".", "_")
	return strings.Title(id)
}

func findSwagCmd() (string, error) { return locateDir("cmd/go-swagger", 5) }
func findSwagPkg() (string, error) { return locateDir("pkg/serveswagger", 5) }
func locateDir(dir string, maxTries int) (string, error) {
	tries := 0
	path := dir
	var err error
	var fi os.FileInfo
	for tries < maxTries {
		fi, err = os.Stat(path)
		if err == nil && fi.IsDir() {
			return filepath.Clean(path), nil
		}
		path = "../" + path
		tries++
	}
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		return "", errors.New("not a dir")
	}
	return "", fmt.Errorf("cannot find %s after %d tries", dir, tries)
}
