package xss

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/utils"
)

var templates = template.Must(template.ParseFiles(
	"./views/partials/safeButtons.gohtml",
	"./views/pages/xss.gohtml",
	"./views/partials/ruleInfo.gohtml",
))

func queryHandler(w http.ResponseWriter, r *http.Request, safety string) (template.HTML, bool) {
	s := r.URL.Query().Get("input")
	if safety == "safe" {
		s = url.QueryEscape(s)
	} else if safety == "noop" {
		return template.HTML("NOOP"), false
	}
	//execute input script
	return template.HTML(s), false

}

func paramsHandler(w http.ResponseWriter, r *http.Request, safety string) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	var s string
	s = splitURL[4] + "/" + splitURL[5]
	if safety == "safe" {
		s = url.QueryEscape(s)
	} else if safety == "noop" {
		return template.HTML("NOOP"), false
	}
	return template.HTML(s), false

}

func xssTemplate(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	var buf bytes.Buffer

	err := templates.ExecuteTemplate(&buf, "xss", pd.Rulebar[pd.Name])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true
}

// Handler is the API handler for XSS
func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	switch splitURL[2] {
	case "query":
		return queryHandler(w, r, splitURL[4])
	case "params":
		return paramsHandler(w, r, splitURL[6])
	default:
		return xssTemplate(w, r, pd)
	}
}
