package main

import (
	"html/template"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"bytes"
	"time"
	unvalidatedRedirect "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/unvalidatedRedirect"
	utils "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/utils"
	ssrf  "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/ssrf"
	nosql  "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/nosqlInjection"
	commandInjection "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/commandInjection"
	sqlInjection "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/sqlInjection"
	pathTraversal "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/pathTraversal"
	xss "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/xss"
)

var pd utils.Parameters
var t *template.Template

const Port = 8080

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
		pd.Body = template.HTML(buf.String())//"<h1> Under Construction</h1>\n <img src=\"https://golang.org/doc/gopher/pencil/gopherswrench.jpg\">"
	}
	pd.Year = time.Now().Year()
	err := t.ExecuteTemplate(w, "layout.gohtml", &pd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, utils.Parameters) (template.HTML, bool), name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var use_layout bool
		pd.Year = time.Now().Year()
		pd.Name = name
		pd.Body, use_layout = fn(w, r, pd)
		if use_layout {
			err := t.ExecuteTemplate(w, "layout.gohtml", &pd)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			fmt.Fprintf(w, string(pd.Body))
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
	var err error
	pd.Port = fmt.Sprintf(":%d", Port)
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
	jsonFile.Close()
	json.Unmarshal([]byte(byteValue), &pd.Rulebar)

	log.Println("Server Startup at: localhost" + pd.Port)
	nosql.MongoInit()
	defer nosql.MongoKill()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ssrf/", makeHandler(ssrf.Handler, "ssrf"))
	http.HandleFunc("/unvalidatedRedirect/", makeHandler(unvalidatedRedirect.Handler, "unvalidatedRedirect"))
	http.HandleFunc("/cmdInjection/", makeHandler(commandInjection.Handler, "cmdInjection"))
	http.HandleFunc("/nosqlInjection/", makeHandler(nosql.Handler, "nosqlInjection"))
	http.HandleFunc("/pathTraversal/", makeHandler(pathTraversal.Handler, "pathTraversal"))
	http.HandleFunc("/sqlInjection/", makeHandler(sqlInjection.Handler, "sqlInjection"))
	http.HandleFunc("/xss/", makeHandler(xss.Handler, "xss"))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public"))))

	log.Fatal(http.ListenAndServe(pd.Port, nil))
}
