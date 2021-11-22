package unvalidated

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

func httpRedirectHandler(w http.ResponseWriter, r *http.Request, pd common.Parameters, splitURL []string) (template.HTML, bool) {
	mode := splitURL[len(splitURL)-1]

	switch mode {
	case "safe":
		formValue := common.GetUserInput(r)
		sanitizedURL := url.PathEscape(formValue)
		http.Redirect(w, r, sanitizedURL, http.StatusFound)
	case "unsafe":
		formValue := common.GetUserInput(r)
		http.Redirect(w, r, formValue, http.StatusFound)
	case "noop":
		http.Redirect(w, r, "http://www.example.com", http.StatusFound)
	}

	return "", false
}

func unvalidatedTemplate(w http.ResponseWriter, r *http.Request, pd common.Parameters) (template.HTML, bool) {
	return "unvalidatedRedirect.gohtml", true
}

// Handler is the API handler for unvalidated redirect
func Handler(w http.ResponseWriter, r *http.Request, pd common.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")

	switch splitURL[2] {

	case "http.Redirect":
		return httpRedirectHandler(w, r, pd, splitURL)

	default:
		return unvalidatedTemplate(w, r, pd)
	}
}
