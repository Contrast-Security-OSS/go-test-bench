package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	Sinks    []Sink   // one per vulnerable function
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

// UnsafeRequests generates an unsafe request for each input and sink defined for this endpoint.
func (r *Route) UnsafeRequests(addr string) ([]*http.Request, error) {
	reqs := make([]*http.Request, 0, len(r.Inputs)*len(r.Sinks))
	for _, s := range r.Sinks {
		if len(s.Name) == 0 || s.Name == "_" {
			continue
		}
		for _, i := range r.Inputs {
			method := methodFromInput(i)
			var u string
			if r.genericTmpl {
				// different parm order, to more easily work with gin
				u = fmt.Sprintf("http://%s%s/%s/%s/unsafe", addr, r.Base, s.Name, i)
			} else {
				u = fmt.Sprintf("http://%s%s/%s/%s/unsafe", addr, r.Base, i, s.Name)
			}
			req, err := http.NewRequest(method, u, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}
			s.AddPayloadToRequest(req, i, "", r.Payload)
			reqs = append(reqs, req)
		}
	}
	return reqs, nil
}

func methodFromInput(in string) string {
	for _, i := range []string{"cookie", "body"} {
		if strings.Contains(in, i) {
			return http.MethodPost
		}
	}
	return http.MethodGet
}

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

// FindViewsDir looks for views dir in working dir or two dirs up, where it's likely to be found in tests
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

var rmap RouteMap

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

type Safety string

const (
	Unsafe Safety = "unsafe"
	Safe   Safety = "safe"
	NOOP   Safety = "noop"
)
