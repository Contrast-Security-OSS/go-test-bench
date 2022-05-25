package common

import (
	"encoding/json"
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

// FindViewsDir looks for views dir in working dir or two dirs up, where it's
// likely to be found in tests.
func FindViewsDir() (string, error) {
	path := "views"
	fi, err := os.Stat(path)
	if err != nil || !fi.IsDir() {
		path = "../../" + path
		fi, err = os.Stat(path)
	}
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		return "", errors.New("not a dir")
	}
	return filepath.Clean(path), nil
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
	rmap = make(RouteMap)
	//add legacy routes
	log.Println("Loading routes.json from ./views/routes.json")
	path, err := FindViewsDir()
	if err != nil {
		log.Fatal("cannot find routes.json")
	}
	path += "/routes.json"
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	if err = json.Unmarshal(data, &rmap); err != nil {
		log.Fatal(err)
	}
	//add migrated routes
	for _, r := range routes {
		if _, ok := rmap[r.Base]; ok {
			log.Println("overwriting route - can be removed from json: ", r.Base, r)
			//or is duplicated within routes
		}
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
