package xss

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	utils "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/utils"
)

var templates = template.Must(template.ParseFiles("./views/partials/safeButtons.gohtml", "./views/pages/xss.gohtml", "./views/partials/ruleInfo.gohtml"))

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
	splitUrl := strings.Split(r.URL.Path, "/")
	var s string
	s = splitUrl[4] + "/" + splitUrl[5]
	if safety == "safe" {
		s = url.QueryEscape(s)
	} else if safety == "noop" {
		return template.HTML("NOOP"), false
	}
	return template.HTML(s), false

}

func defaultHandler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "xss", pd.Rulebar[pd.Name])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true
}

func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	splitUrl := strings.Split(r.URL.Path, "/")
	switch splitUrl[2] {
	case "query":
		return queryHandler(w, r, splitUrl[4])
	case "params":
		return paramsHandler(w, r, splitUrl[6])
	default:
		return defaultHandler(w, r, pd)
	}
}
