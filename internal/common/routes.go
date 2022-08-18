package common

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Verbose increases the verbosity of logging.
var Verbose bool

// Route is the template information for a specific route
type Route struct {
	Name     string   // human-readable name
	Link     string   // owasp link
	Base     string   // short name, suitable for use in filename or URL - i.e. cmdInjection
	TmplFile string   // name of template used for non-result page; default is Base + '.gohtml'
	Products []string // relevant Contrast products
	Inputs   []string // input methods supported by this app: query, cookies, body, headers, headers-json, ...
	Sinks    []*Sink  // one per vulnerable function
	Payload  string   // must be set for the default template.

	genericTmpl bool
}

func (r *Route) String() string {
	lines := []string{fmt.Sprintf("%s %q %s", r.Base, r.Name, r.Link)}
	for _, s := range r.Sinks {
		if str := s.String(); len(str) > 0 {
			lines = append(lines, "- "+str)
		}
	}
	return strings.Join(lines, "    \n")
}

// UsesGenericTmpl returns true if the route uses the generic vulnerability template.
func (r *Route) UsesGenericTmpl() bool { return r.genericTmpl }

// RouteMap is a map from base path to Route
type RouteMap map[string]Route

// Routes is a slice of Route
type Routes []Route

// sortable for swagger code gen
func (rs Routes) Len() int           { return len(rs) }
func (rs Routes) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs Routes) Less(i, j int) bool { return rs[i].Name < rs[j].Name }

func (rs Routes) String() string {
	var list []string
	for _, r := range rs {
		list = append(list, r.String())
	}
	return strings.Join(list, "\n")
}

// AllRoutes contains all "new" (not in json) routes.
var AllRoutes Routes

// Register adds one or more Endpoints to the global list of routes.
func Register(r Route) {
	if len(r.Base) == 0 || len(r.Name) == 0 {
		log.Fatalf("Base and Name must both be populated in %v", r)
	}
	if strings.Contains(r.Base, "/") {
		log.Fatalf("%s: slashes not allowed in Base", r.Name)
	}
	if len(r.TmplFile) == 0 {
		templatesDir, err := FindViewsDir()
		if err != nil {
			log.Fatal("cannot find views dir:", err)
		}
		p := filepath.Join(templatesDir, "pages", r.Base+".gohtml")
		if _, err := os.Stat(p); err != nil {
			//does not exist - use generic
			r.genericTmpl = true
			r.TmplFile = "rule.gohtml"
		} else {
			r.TmplFile = r.Base + ".gohtml"
		}
	}
	r.Base = "/" + r.Base
	for i, s := range r.Sinks {
		if (s.Handler == nil) == (s.VulnerableFnWrapper == nil) {
			log.Fatalf("sink #%d in %#v: exactly one of {Handler, VulnerableFnWrapper} must be set", i, r)
		}
		if len(s.Name) == 0 {
			log.Fatalf("0-len sink name at %d in %#v", i, r)
		}
		if strings.Contains(s.URL, "/") {
			log.Fatal("slashes not allowed in sink url:", s, r)
		}
		if len(s.URL) == 0 {
			r.Sinks[i].URL = s.Name
		}
	}
	AllRoutes = append(AllRoutes, r)
}

// FindViewsDir looks for the views dir, which contains our html templates.
// It looks in the current dir and its parents.
func FindViewsDir() (string, error) { return LocateDir("views", 5) }

// LocateDir finds a dir with the given name and returns its path.
// The given name may contain a slash, i.e. 'cmd/go-swagger'.
func LocateDir(dir string, maxTries int) (string, error) {
	tries := 0
	path := dir
	var err error
	var fi os.FileInfo
	for tries < maxTries {
		fi, err = os.Stat(path)
		if err == nil && fi.IsDir() {
			return filepath.Clean(path), nil
		}
		path = filepath.Join("..", path)
		tries++
	}
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		return "", errors.New("not a dir")
	}
	return "", fmt.Errorf("cannot find %s in any of %d parent dirs", dir, tries)
}

// Templates is the map we use to lookup the parsed templates
// based on filenames. It is intended for use for use
// by all frameworks supported by the bench.
var Templates = make(map[string]*template.Template)

// ParseViewTemplates is used to set up the template
// resources for use by std and go-swagger
func ParseViewTemplates() error {
	templatesDir, err := FindViewsDir()
	if err != nil {
		return err
	}

	pages, err := filepath.Glob(filepath.Join(templatesDir, "pages", "*.gohtml"))
	if err != nil {
		return err
	}
	if len(pages) == 0 {
		log.Fatal("nothing found in ./views/pages")
	}
	partials, err := filepath.Glob(filepath.Join(templatesDir, "partials", "*.gohtml"))
	if err != nil {
		return err
	}
	if len(partials) == 0 {
		log.Fatal("nothing found in ./views/partials")
	}
	layout := filepath.Join(templatesDir, "layout.gohtml")

	fmap := FuncMap()

	for _, p := range pages {
		files := append([]string{layout, p}, partials...)
		tmpl, err := template.New(p).Funcs(fmap).ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}
		Templates[filepath.Base(p)] = tmpl
	}

	return nil
}

// Reset clears AllRoutes and rmap. For testing.
func Reset() {
	AllRoutes = nil
	rmap = nil
}

var rmap RouteMap

// GetRouteMap returns the already-populated RouteMap.
func GetRouteMap() RouteMap {
	return rmap
}

// PopulateRouteMap returns a RouteMap, for use in nav bar template.
func PopulateRouteMap(routes Routes) RouteMap {
	rmap = make(RouteMap, len(routes))
	for _, r := range routes {
		rmap[r.Base] = r
	}
	//clean up
	var rts []string
	for k := range rmap {
		rts = append(rts, k)
	}
	for _, k := range rts { //don't range over rmap directly, as we need to be able to delete from it
		r := rmap[k]
		for i := 0; i < len(r.Products); i++ {
			if r.Products[i] == "Protect" {
				//remove - product not yet available
				r.Products = append(r.Products[:i], r.Products[i+1:]...)
				continue
			}
		}
		if len(r.Products) == 0 {
			//don't show anything that was protect-only
			delete(rmap, k)
		} else {
			rmap[k] = r
		}
	}
	// print route list?
	if Verbose {
		var lines []string
		for _, r := range rmap {
			if len(r.Sinks) > 0 && len(r.Sinks[0].Name) > 0 {
				lines = append(lines, r.String())
			}
		}
		log.Printf("vulnerable routes:\n%s", strings.Join(lines, "\n"))
	}
	return rmap
}

// Safety indicates whether input to the vulnerable function will be sanitized
// or not, or if the vulnerable func will be bypassed entirely.
type Safety string

const (
	// Unsafe indicates no sanitization will be performed.
	Unsafe Safety = "unsafe"
	// Safe indicates input will be sanitized.
	Safe Safety = "safe"
	// NOOP indicates the vulnerable function will not be called.
	NOOP Safety = "noop"
)
