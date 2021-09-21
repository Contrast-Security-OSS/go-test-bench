package main

import (
	"html/template"

	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Contrast-Security-OSS/go-test-bench/cmdi"
	"github.com/Contrast-Security-OSS/go-test-bench/pathtraversal"
	"github.com/Contrast-Security-OSS/go-test-bench/sqli"
	"github.com/Contrast-Security-OSS/go-test-bench/ssrf"
	"github.com/Contrast-Security-OSS/go-test-bench/unvalidated"
	"github.com/Contrast-Security-OSS/go-test-bench/utils"
	"github.com/Contrast-Security-OSS/go-test-bench/xss"
)

var pd = utils.Parameters{
	Year: 2020,
	Logo: "https://blog.golang.org/gopher/header.jpg",
}
var templates = make(map[string]*template.Template)

// DefaultPort is the port that the API runs on if no command line argument is specified
const DefaultPort = 8080

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	if r.URL.Path == "/" {
		t = templates["index.gohtml"]
	} else {
		t = templates["underConstruction.gohtml"]
	}
	err := t.ExecuteTemplate(w, "layout.gohtml", pd)
	if err != nil {
		log.Print(err.Error())
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, utils.Parameters) (template.HTML, bool), name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pd.Name = name
		data, useLayout := fn(w, r, pd)
		if useLayout {
			err := templates[string(data)].ExecuteTemplate(w, "layout.gohtml", &pd)
			if err != nil {
				log.Print(err.Error())
			}
		} else {
			fmt.Fprint(w, data)
		}
	}
}

func parseTemplates() error {
	templatesDir := filepath.Clean("../../views")
	pages, err := filepath.Glob(filepath.Join(templatesDir, "pages", "*.gohtml"))
	if err != nil {
		return err
	}
	partials, err := filepath.Glob(filepath.Join(templatesDir, "partials", "*.gohtml"))
	if err != nil {
		return err
	}
	layout := filepath.Join(templatesDir, "layout.gohtml")

	for _, p := range pages {
		files := append([]string{layout, p}, partials...)
		templates[filepath.Base(p)] = template.Must(template.ParseFiles(files...))
	}

	return nil
}

func main() {

	// Setup command line flags
	portPtr := flag.Int("port", DefaultPort, "listen on this port")
	flag.Parse()
	port := *portPtr

	pd.Port = fmt.Sprintf(":%d", port)
	log.Println("Loading templates...")
	err := parseTemplates()
	if err != nil {
		log.Fatalln("Cannot parse templates:", err)
	}
	log.Println("Templates loaded.")

	log.Println("Loading routes.json from /views/routes.json")
	jsonFile, err := os.Open("../../views/routes.json")
	if err != nil {
		log.Fatalln(err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	_ = jsonFile.Close()
	_ = json.Unmarshal([]byte(byteValue), &pd.Rulebar)

	log.Println("Server startup at: localhost" + pd.Port)

	// Attempt to connect to MongoDB with a 30 second timeout
	// err = nosql.MongoInit(time.Second * 30)
	// if err != nil {
	// 	log.Printf("Could not connect the Mongo client: err = %s", err)
	// 	os.Exit(1)
	// }

	//defer nosql.MongoKill()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ssrf/", makeHandler(ssrf.Handler, "ssrf"))
	http.HandleFunc("/unvalidatedRedirect/", makeHandler(unvalidated.Handler, "unvalidatedRedirect"))
	http.HandleFunc("/cmdInjection/", makeHandler(cmdi.Handler, "cmdInjection"))

	// http.HandleFunc("/nosqlInjection/", makeHandler(nosql.Handler, "nosqlInjection"))

	http.HandleFunc("/pathTraversal/", makeHandler(pathtraversal.Handler, "pathTraversal"))
	http.HandleFunc("/sqlInjection/", makeHandler(sqli.Handler, "sqlInjection"))
	http.HandleFunc("/xss/", makeHandler(xss.Handler, "xss"))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../../public"))))

	log.Fatal(http.ListenAndServe(pd.Port, nil))
}
