package servestd

import (
	"html/template"
	"strings"

	"fmt"
	"log"
	"net/http"
	"path/filepath"

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
	Year:      2021,
	Logo:      "https://blog.golang.org/gopher/header.jpg",
	Framework: "stdlib",
}

var templates = make(map[string]*template.Template)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	if r.URL.Path == "/" {
		t = templates["index.gohtml"]
	} else {
		t = templates["underConstruction.gohtml"]
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
			err := templates[string(data)].ExecuteTemplate(w, "layout.gohtml", &parms)
			if err != nil {
				log.Print(err.Error())
			}
		} else {
			fmt.Fprint(w, data)
		}
	}
}

func newHandler(v common.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(v.Name, r.URL.Path)
		var parms = common.Parameters{
			ConstParams: Pd,
			Name:        v.Base,
		}
		var data = template.HTML(v.TmplFile)
		isTmpl := true
		elems := strings.Split(r.URL.Path, "/")
		// To figure out whether we're serving a sink or the main page, check the
		// element with index 2 against each Sink.URL; if no match, serve main page.
		// Seems like there should be a less ugly way...
		for _, s := range v.Sinks {
			if len(elems) > 2 && elems[2] == s.URL {
				mode := elems[len(elems)-1]
				data, isTmpl = s.Handler(mode, common.GetUserInput(r))
				break
			}
		}
		if isTmpl {
			err := templates[string(data)].ExecuteTemplate(w, "layout.gohtml", &parms)
			if err != nil {
				log.Print(err.Error())
			}
		} else {
			fmt.Fprint(w, data)
		}
	}
}

func parseTemplates() error {
	templatesDir, err := common.FindViewsDir()
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

	fmap := template.FuncMap{"tolower": strings.ToLower}

	for _, p := range pages {
		files := append([]string{layout, p}, partials...)
		tmpl, err := template.New(p).Funcs(fmap).ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}
		templates[filepath.Base(p)] = tmpl
	}

	return nil
}

// Setup loads templates, sets up routes, etc.
func Setup() {
	log.Println("Loading templates...")
	err := parseTemplates()
	if err != nil {
		log.Fatalln("Cannot parse templates:", err)
	}
	log.Println("Templates loaded.")

	//register all routes at this point.
	cmdi.RegisterRoutes("stdlib")

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

	http.HandleFunc("/pathTraversal/", makeHandler(pathtraversal.Handler, "pathTraversal"))
	http.HandleFunc("/sqlInjection/", makeHandler(sqli.Handler, "sqlInjection"))
	http.HandleFunc("/xss/", makeHandler(xss.Handler, "xss"))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./public"))))
}
