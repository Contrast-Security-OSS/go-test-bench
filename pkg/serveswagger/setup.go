package serveswagger

import (
	"log"
	"net/http"
	"strings"
	"path/filepath"
	"html/template"
	"github.com/go-openapi/runtime"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// This lets us implement the Responder interface from go-swagger
// so we can have access to the http.ResponseWriter when each route gets exercised.
type CustomResponder func(http.ResponseWriter, runtime.Producer)

func (c CustomResponder) WriteResponse(w http.ResponseWriter, p runtime.Producer) {
	c(w, p)
}

var Pd = common.ConstParams{
	Year:      	2022,
	Logo: 		"https://raw.githubusercontent.com/swaggo/swag/master/assets/swaggo.png",
	Framework: 	"Go-Swagger",
}

var Templates = make(map[string]*template.Template)

// took the parseTemplates procedure from std-lib. Unused for this moment.
// Meant to serve as a starting point when setting up templates for swagger.
// Right now relying on directly on std for templates
func ParseTemplates() error {
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
		Templates[filepath.Base(p)] = tmpl
	}

	return nil
}
