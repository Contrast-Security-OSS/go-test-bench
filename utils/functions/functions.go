package functions

import (
	"html/template"
	"net/http"
	"time"
)

type Route struct {
	Base string
	Name string
	Link string
	Products []string
	Inputs []string
	Sinks string
}

type Params struct {
	Body string
	Year int
	Rulebar map[string]Route
}


var T *template.Template

func RenderTemplate(w http.ResponseWriter, tmpl string, p *Params) {
	p.Year = time.Now().Year()
	p.Body = "Hello world!" //probably should not go here
	err := T.ExecuteTemplate(w, tmpl+".gohtml", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}