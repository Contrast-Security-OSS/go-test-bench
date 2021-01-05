package main

import (
	"html/template"

	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Contrast-Security-OSS/go-test-bench/cmdi"
	"github.com/Contrast-Security-OSS/go-test-bench/nosql"
	"github.com/Contrast-Security-OSS/go-test-bench/pathtraversal"
	"github.com/Contrast-Security-OSS/go-test-bench/sqli"
	"github.com/Contrast-Security-OSS/go-test-bench/ssrf"
	"github.com/Contrast-Security-OSS/go-test-bench/unvalidated"
	"github.com/Contrast-Security-OSS/go-test-bench/utils"
	"github.com/Contrast-Security-OSS/go-test-bench/xss"
)

var pd utils.Parameters
var t *template.Template

// DefaultPort is the port that the API runs on if no command line argument is specified
const DefaultPort = 8080

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		pd.Body = "<h1 class=\"page-header\">Contrast Go Test Bench</h1>\n <img src=\"https://blog.golang.org/gopher/header.jpg\">"
	} else {
		//var templates = template.Must(template.ParseFiles("./views/pages/underConstruction.gohtml"))
		var buf bytes.Buffer
		err := t.ExecuteTemplate(&buf, "underConstruction.gohtml", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		pd.Body = template.HTML(buf.String()) //"<h1> Under Construction</h1>\n <img src=\"https://golang.org/doc/gopher/pencil/gopherswrench.jpg\">"
	}
	pd.Year = time.Now().Year()
	err := t.ExecuteTemplate(w, "layout.gohtml", &pd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, utils.Parameters) (template.HTML, bool), name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var useLayout bool
		pd.Year = time.Now().Year()
		pd.Name = name
		pd.Body, useLayout = fn(w, r, pd)
		if useLayout {
			err := t.ExecuteTemplate(w, "layout.gohtml", &pd)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			_, _ = fmt.Fprint(w, string(pd.Body))
		}

	}
}

func parseTemplates() (*template.Template, error) {
	templ := template.New("")
	err := filepath.Walk("./views", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".gohtml") {
			if _, err = templ.ParseFiles(path); err != nil {
				log.Println(err)
			}
			log.Println("Loading - " + path)
		}
		return err
	})
	return templ, err
}

func main() {

	// Setup command line flags
	portPtr := flag.Int("port", DefaultPort, "listen on this port")
	flag.Parse()
	port := *portPtr

	var err error
	pd.Port = fmt.Sprintf(":%d", port)
	log.Println("Loading Templates: ")
	log.Println("----------------------------------")
	t, err = parseTemplates()
	if err != nil {
		log.Fatalln("Cannot parse templates:", err)
	}
	log.Println("----------------------------------")
	log.Println("Templates Loaded ")

	log.Println("Loading routes.json from /views/routes.json")
	jsonFile, err := os.Open("./views/routes.json")
	if err != nil {
		log.Fatalln(err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	_ = jsonFile.Close()
	_ = json.Unmarshal([]byte(byteValue), &pd.Rulebar)

	log.Println("Server Startup at: localhost" + pd.Port)
	nosql.MongoInit()
	defer nosql.MongoKill()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ssrf/", makeHandler(ssrf.Handler, "ssrf"))
	http.HandleFunc("/unvalidatedRedirect/", makeHandler(unvalidated.Handler, "unvalidatedRedirect"))
	http.HandleFunc("/cmdInjection/", makeHandler(cmdi.Handler, "cmdInjection"))
	http.HandleFunc("/nosqlInjection/", makeHandler(nosql.Handler, "nosqlInjection"))
	http.HandleFunc("/pathTraversal/", makeHandler(pathtraversal.Handler, "pathTraversal"))
	http.HandleFunc("/sqlInjection/", makeHandler(sqli.Handler, "sqlInjection"))
	http.HandleFunc("/xss/", makeHandler(xss.Handler, "xss"))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public"))))

	log.Fatal(http.ListenAndServe(pd.Port, nil))
}
