package servejschmidt

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
	"github.com/julienschmidt/httprouter"
)

// Pd is unchanging parameter data shared between all routes.
var Pd = common.ConstParams{
	Year:      2022,
	Logo:      "https://blog.golang.org/gopher/header.jpg",
	Framework: "Julienschmidt",
}

func newHandler(v common.Route) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Println(v.Name, r.URL.Path)

		// Split off the first element from elems, which should be the mode.
		mode := strings.Split(strings.Trim(p.ByName("elems"), "/"), "/")[0]
		for _, s := range v.Sinks {
			if p.ByName("sink") != s.URL {
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

func add(router *httprouter.Router, rt common.Route) {
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
	router.GET(rt.Base+"/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	// Julienshmidt only allows one request method for each route,
	// so each input method has to be a separate route.
	postInputs := map[string]struct{}{ // These input types use a POST request instead of GET
		"body":          {},
		"buffered-body": {},
		"cookies":       {},
	}
	for _, input := range rt.Inputs {
		if _, ok := postInputs[input]; ok {
			router.POST(rt.Base+"/:sink/"+input+"/*elems", newHandler(rt))
		} else {
			router.GET(rt.Base+"/:sink/"+input+"/*elems", newHandler(rt))
		}
	}

}

type httpHandlerPair struct {
	http.ResponseWriter
	*http.Request
}

// RegisterRoutes registers all decoupled routes used by servejschmidt. Shared with cmd/exercise.
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
func Setup() *httprouter.Router {
	log.Println("Loading Templates...")
	err := common.ParseViewTemplates()
	if err != nil {
		log.Fatalln("Cannot parse Templates:", err)
	}
	log.Println("Templates loaded.")

	// register all routes in this function
	RegisterRoutes()

	Pd.Rulebar = common.PopulateRouteMap(common.AllRoutes)

	router := httprouter.New()

	log.Println("Server startup at: " + Pd.Addr)

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	router.ServeFiles("/assets/*filepath", http.Dir(pub))

	return router

}
