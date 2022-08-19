package servestd

import (
	"html/template"
	"strings"

	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/sqli"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/pathtraversal"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/ssrf"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/unvalidated"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/xss"
)

//Pd is unchanging parameter data shared between all routes.
var Pd = common.ConstParams{
	Year:      2022,
	Logo:      "https://blog.golang.org/gopher/header.jpg",
	Framework: "stdlib",
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	if r.URL.Path == "/" {
		t = common.Templates["index.gohtml"]
	} else {
		t = common.Templates["pageUnsupported.gohtml"]
		w.WriteHeader(http.StatusNotFound)
	}
	w.Header().Set("Application-Framework", "Stdlib")
	err := t.ExecuteTemplate(w, "layout.gohtml", Pd)
	if err != nil {
		log.Print(err.Error())
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, common.Parameters) (template.HTML, bool), name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var parms = common.Parameters{
			ConstParams: Pd,
			Name:        name,
		}
		data, useLayout := fn(w, r, parms)
		if useLayout {
			err := common.Templates[string(data)].ExecuteTemplate(w, "layout.gohtml", &parms)
			if err != nil {
				log.Print(err.Error())
			}
		} else {
			fmt.Fprint(w, data)
		}
	}
}

func newHandler(v common.Route) http.HandlerFunc {
	for _, s := range v.Sinks {
		if s.Handler == nil {
			var err error
			s.Handler, err = common.GenericHandler(s)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(v.Name, r.URL.Path)
		var parms = common.Parameters{
			ConstParams: Pd,
			Name:        v.Base,
		}
		elems := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(elems) < 2 {
			// main page
			err := common.Templates[v.TmplFile].ExecuteTemplate(w, "layout.gohtml", &parms)
			if err != nil {
				log.Print(err.Error())
				fmt.Fprintf(w, "template error: %s", err)
			}
			return
		}
		for _, s := range v.Sinks {
			if elems[1] != s.URL {
				continue
			}
			// can't assume safety is last if input is via path parameters.
			safeIdx := 3
			if len(elems) <= safeIdx {
				safeIdx = len(elems) - 1
			}
			mode := common.Safety(elems[safeIdx])
			switch mode {
			case common.NOOP, common.Safe, common.Unsafe:
				// valid modes
			default:
				// invalid
				w.WriteHeader(http.StatusNotFound)
				return
			}

			in := common.GetUserInput(r)
			data, mime, status := s.Handler(mode, in, httpHandlerPair{w, r})
			if len(data) == 0 {
				// don't unconditionally write response, as it can result in
				// - an error log of "http: superfluous response.WriteHeader call"
				return
			}
			w.WriteHeader(status)
			if len(mime) == 0 {
				mime = "text/plain"
			}
			w.Header().Set("Content-Type", mime)
			w.Header().Set("Cache-Control", "no-store") //makes development a whole lot easier
			fmt.Fprint(w, data)
			return
		}
		// does not match any sink or the main page
		w.WriteHeader(http.StatusNotFound)
	}
}

type httpHandlerPair struct {
	http.ResponseWriter
	*http.Request
}

// RegisterRoutes registers all decoupled routes used by servestd. Shared with cmd/exercise.
func RegisterRoutes() {
	cmdi.RegisterRoutes()
	sqli.RegisterRoutes()
	pathtraversal.RegisterRoutes()
	ssrf.RegisterRoutes()
	unvalidated.RegisterRoutes(&common.Sink{
		Name:     "http.Redirect",
		Sanitize: url.PathEscape,
		VulnerableFnWrapper: func(opaque interface{}, payload string) (data string, raw bool, err error) {
			p, ok := opaque.(httpHandlerPair)
			if !ok {
				log.Fatalf("'opaque': want httpHandlerPair, got %T", opaque)
			}
			w, r := p.ResponseWriter, p.Request
			http.Redirect(w, r, payload, http.StatusFound)
			return "", true, nil
		},
	})
	xss.RegisterRoutes()
}

// Setup loads templates, sets up routes, etc.
func Setup() {
	log.Println("Loading Templates...")
	err := common.ParseViewTemplates()
	if err != nil {
		log.Fatalln("Cannot parse Templates:", err)
	}
	log.Println("Templates loaded.")

	// register all routes in this function
	RegisterRoutes()

	Pd.Rulebar = common.PopulateRouteMap(common.AllRoutes)

	log.Println("Server startup at: " + Pd.Addr)

	http.HandleFunc("/", rootHandler)

	for _, r := range common.AllRoutes {
		http.HandleFunc(r.Base+"/", newHandler(r))
	}

	pub, err := common.LocateDir("public", 5)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(pub))))
}
