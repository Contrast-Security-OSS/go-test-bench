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

var Verbose bool

// HandlerFn is a framework-agnostic function to handle a vulnerable endpoint.
type HandlerFn func(mode, in string) (template.HTML, bool)

// Sink is a struct that identifies the name
// of the sink, the associated URL and the
// HTTP method
type Sink struct {
	Name    string
	URL     string
	Method  string
	Handler HandlerFn // the vulnerable function which recieves unsanitized input
}

func (s *Sink) String() string {
	return fmt.Sprintf("%s: %s .../%s", s.Name, s.Method, s.URL)
}

// Route is the template information for a specific route
type Route struct {
	Name     string   // human-readable name
	Link     string   // owasp link
	Base     string   // short name, suitable for use in filename or URL - i.e. cmdInjection
	TmplFile string   // name of template used for non-result page; default is Base + '.gohtml'
	Products []string // relevant Contrast products
	Inputs   []string // input methods supported by this app
	Sinks    []Sink   // one per vulnerable function
}

func (r *Route) String() string {
	lines := []string{fmt.Sprintf("%q %s %s", r.Name, r.Base, r.Link)}
	for _, s := range r.Sinks {
		if str := s.String(); len(str) > 0 {
			lines = append(lines, "- "+str)
		}
	}
	return strings.Join(lines, "    \n")
}

type RouteMap map[string]Route
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
		r.TmplFile = r.Base + ".gohtml"
	}
	r.Base = "/" + r.Base
	for i, s := range r.Sinks {
		if len(r.Name) == 0 {
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

//Look for views dir in working dir
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

// PopulateRouteMap returns a RouteMap, for use in nav bar template.
func PopulateRouteMap(routes Routes) (rmap RouteMap) {
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
	var lines []string
	for _, r := range rmap {
		if len(r.Sinks) > 0 && len(r.Sinks[0].Name) > 0 {
			lines = append(lines, r.String())
		}
	}
	if Verbose {
		log.Printf("vulnerable routes:\n%s", strings.Join(lines, "\n"))
	}
	return
}
