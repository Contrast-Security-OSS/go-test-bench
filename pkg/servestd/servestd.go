package servestd

import (
	"html/template"
	"strings"

	"fmt"
	"log"
	"net/http"

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
		t = common.Templates["underConstruction.gohtml"]
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
			if elems[1] == s.URL {
				mode := common.Safety(elems[len(elems)-1])
				switch mode {
				case common.NOOP, common.Safe, common.Unsafe:
					// valid modes
				default:
					// invalid
					w.WriteHeader(http.StatusNotFound)
					return
				}

				in := common.GetUserInput(r)
				data, status := s.Handler(mode, in, nil)
				w.WriteHeader(status)
				w.Header().Set("Cache-Control", "no-store") //makes development a whole lot easier
				fmt.Fprint(w, data)
				return
			}
		}
		// does not match any sink or the main page
		w.WriteHeader(http.StatusNotFound)
	}
}

// Setup loads templates, sets up routes, etc.
func Setup() {
	log.Println("Loading Templates...")
	err := common.ParseViewTemplates()
	if err != nil {
		log.Fatalln("Cannot parse Templates:", err)
	}
	log.Println("Templates loaded.")

	// register all routes at this point.
	cmdi.RegisterRoutes()
	sqli.RegisterRoutes()
	pathtraversal.RegisterRoutes()

	Pd.Rulebar = common.PopulateRouteMap(common.AllRoutes)

	log.Println("Server startup at: " + Pd.Addr)

	// Attempt to connect to MongoDB with a 30 second timeout
	// err = nosql.MongoInit(time.Second * 30)
	// if err != nil {
	// 	log.Printf("Could not connect the Mongo client: err = %s", err)
	// 	os.Exit(1)
	// }

	//defer nosql.MongoKill()

	http.HandleFunc("/", rootHandler)

	for _, r := range common.AllRoutes {
		http.HandleFunc(r.Base+"/", newHandler(r))
	}

	http.HandleFunc("/ssrf/", makeHandler(ssrf.Handler, "ssrf"))
	http.HandleFunc("/unvalidatedRedirect/", makeHandler(unvalidated.Handler, "unvalidatedRedirect"))

	// http.HandleFunc("/nosqlInjection/", makeHandler(nosql.Handler, "nosqlInjection"))

	http.HandleFunc("/xss/", makeHandler(xss.Handler, "xss"))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./public"))))
}
