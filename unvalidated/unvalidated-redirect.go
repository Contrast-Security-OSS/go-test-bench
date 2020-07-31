package unvalidated

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
	"./views/pages/unvalidatedRedirect.gohtml",
	"./views/partials/ruleInfo.gohtml",
))

func httpRedirectHandler(w http.ResponseWriter, r *http.Request, pd utils.Parameters, splitURL []string) (template.HTML, bool) {
	mode := splitURL[len(splitURL)-1]
	formValue := r.FormValue("input")

	if mode == "unsafe" {
		http.Redirect(w, r, formValue, http.StatusFound)

	} else if mode == "noop" {
		http.Redirect(w, r, "http://www.example.com", http.StatusFound)

	} else if mode == "safe" {
		sanatizedURL := url.PathEscape(formValue)
		http.Redirect(w, r, sanatizedURL, http.StatusFound)
	}
	return "", false

}

func unvalidatedTemplate(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	var buf bytes.Buffer

	err := templates.ExecuteTemplate(&buf, "unvalidatedRedirect", pd.Rulebar[pd.Name])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true
}

// Handler is the API handler for unvalidated redirect
func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")

	switch splitURL[2] {

	case "http.Redirect":
		return httpRedirectHandler(w, r, pd, splitURL)

	default:
		return unvalidatedTemplate(w, r, pd)
	}
}
