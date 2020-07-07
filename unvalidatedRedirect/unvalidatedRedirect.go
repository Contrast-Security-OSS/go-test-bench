package unvalidatedRedirect

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	utils "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/utils"
)

var templates = template.Must(template.ParseFiles("./views/partials/safeButtons.gohtml", "./views/pages/unvalidatedRedirect.gohtml", "./views/partials/ruleInfo.gohtml"))

func httpRedirectHandler(w http.ResponseWriter, r *http.Request, pd utils.Parameters, split_url []string) (template.HTML, bool) {
	mode := split_url[len(split_url)-1]
	form_value := r.FormValue("input")

	if mode == "unsafe" {
		http.Redirect(w, r, form_value, http.StatusFound)

	} else if mode == "noop" {
		http.Redirect(w, r, "http://www.example.com", http.StatusFound)

	} else if mode == "safe" {
		sanatized_url := url.PathEscape(form_value)
		http.Redirect(w, r, sanatized_url, http.StatusFound)
	}
	return "", false

}

func defaultHandler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "unvalidatedRedirect", pd.Rulebar[pd.Name])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true
}

func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	split_url := strings.Split(r.URL.Path, "/")

	switch split_url[2] {

	case "http.Redirect":
		return httpRedirectHandler(w, r, pd, split_url)

	default:
		return defaultHandler(w, r, pd)
	}
}
