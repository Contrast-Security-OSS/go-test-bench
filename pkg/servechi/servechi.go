package servechi

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/sqli"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/pathtraversal"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/ssrf"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/unvalidated"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/xss"
	"github.com/go-chi/chi"
)

// Pd is unchanging parameter data shared between all routes.
var Pd = common.ConstParams{
	Year:      2022,
	Logo:      "https://blog.golang.org/gopher/header.jpg",
	Framework: "Chi",
}

func newHandler(v common.Route) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(v.Name, r.URL.Path)

		mode := chi.URLParam(r, "mode")
		for _, s := range v.Sinks {
			if chi.URLParam(r, "sink") != s.URL {
				continue
			}
			mode := common.Safety(mode)
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

func add(router chi.Router, rt common.Route) {
	for _, s := range rt.Sinks {
		if s.Handler == nil {
			var err error
			s.Handler, err = common.GenericHandler(s)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	// main page
	router.Get(rt.Base+"/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(rt.Name, r.URL.Path)
		var parms = common.Parameters{
			ConstParams: Pd,
			Name:        rt.Base,
		}
		err := common.Templates[rt.TmplFile].ExecuteTemplate(w, "layout.gohtml", &parms)
		if err != nil {
			log.Print(err.Error())
			fmt.Fprintf(w, "template error: %s", err)
		}
	})

	postInputs := map[string]struct{}{ // These input types use a POST request instead of GET
		"body":          {},
		"buffered-body": {},
		"cookies":       {},
	}
	for _, input := range rt.Inputs {
		if _, ok := postInputs[input]; ok {
			router.Post(rt.Base+"/{sink}/"+input+"/{mode}", newHandler(rt))
			router.Post(rt.Base+"/{sink}/"+input+"/{mode}/*", newHandler(rt))
		} else {
			router.Get(rt.Base+"/{sink}/"+input+"/{mode}", newHandler(rt))
			router.Get(rt.Base+"/{sink}/"+input+"/{mode}/*", newHandler(rt))
		}
	}

}

type httpHandlerPair struct {
	http.ResponseWriter
	*http.Request
}

// RegisterRoutes registers all decoupled routes used by servechi. Shared with cmd/exercise.
// Right now, this func has the same behavior as servestd.RegisterRoutes
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
func Setup() chi.Router {
	log.Println("Loading Templates...")
	err := common.ParseViewTemplates()
	if err != nil {
		log.Fatalln("Cannot parse Templates:", err)
	}
	log.Println("Templates loaded.")

	// register all routes in this function
	RegisterRoutes()

	Pd.Rulebar = common.PopulateRouteMap(common.AllRoutes)

	router := chi.NewRouter()

	log.Println("Server startup at: " + Pd.Addr)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t := common.Templates["index.gohtml"]
		w.Header().Set("Application-Framework", "Julienschmidt")
		err := t.ExecuteTemplate(w, "layout.gohtml", Pd)
		if err != nil {
			log.Print(err.Error())
		}
	})

	for _, r := range common.AllRoutes {
		add(router, r)
	}

	pub, err := common.LocateDir("public", 5)
	if err != nil {
		log.Fatal(err)
	}
	FileServer(router, "/assets", http.Dir(pub))

	return router

}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// Taken from github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
