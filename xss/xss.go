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
	input := utils.GetUserInput(r)
	if safety == "safe" {
		input = url.QueryEscape(input)
	} else if safety == "noop" {
		return template.HTML("NOOP"), false
	}
	//execute input script
	return template.HTML(input), false

}

// bufferedQueryHandler used as a handler which uses bytes.Buffer for source input ignoring the user input
func bufferedQueryHandler(w http.ResponseWriter, r *http.Request, safety string) (template.HTML, bool) {
	var buf bytes.Buffer
	buf.Write([]byte("<script>"))
	buf.WriteString("alert('buffered input - ")
	buf.WriteRune('©')
	buf.WriteString("');")
	buf.WriteString("</script")
	buf.WriteByte(byte('>'))

	input := string(buf.Bytes())

	if safety == "safe" {
		input = url.QueryEscape(input)
	} else if safety == "noop" {
		return template.HTML("NOOP"), false
	}
	return template.HTML(input), false
}

func paramsHandler(w http.ResponseWriter, r *http.Request, safety string) (template.HTML, bool) {
	// since the attack mode as a last parameter in the query path - e.g. /unsafe, /safe, /noop
	// the user input is placed in the middle and it includes the "/" symbol so we need to combine two pieces
	// /xss/params/reflectedXss/<script>alert(1);</script>/unsafe
	// therefore we specify exact positions of the path to be considered as the user input value
	input := utils.GetPathValue(r, 4, 5)
	if safety == "safe" {
		input = url.QueryEscape(input)
	} else if safety == "noop" {
		return template.HTML("NOOP"), false
	}
	return template.HTML(input), false

}

func bodyHandler(w http.ResponseWriter, r *http.Request, safety string) (template.HTML, bool) {
	if r.Method == http.MethodGet {
		return template.HTML("Cannot GET " + r.URL.Path), false
	}

	input := utils.GetUserInput(r)

	if safety == "safe" {
		input = url.QueryEscape(input)
	} else if safety == "noop" {
		return template.HTML("NOOP"), false
	}

	return template.HTML(input), false
}

// bufferedBodyHandler used as a handler which uses bytes.Buffer for source input
func bufferedBodyHandler(w http.ResponseWriter, r *http.Request, safety string) (template.HTML, bool) {
	if r.Method == http.MethodGet {
		return template.HTML("Cannot GET " + r.URL.Path), false
	}

	buf := bytes.NewBufferString(utils.GetUserInput(r))
	input := buf.String()

	if safety == "safe" {
		input = url.QueryEscape(input)
	} else if safety == "noop" {
		return template.HTML("NOOP"), false
	}

	return template.HTML(input), false
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
	case "buffered-query":
		return bufferedQueryHandler(w, r, splitURL[4])
	case "params":
		return paramsHandler(w, r, splitURL[6])
	case "body":
		return bodyHandler(w, r, splitURL[4])
	case "buffered-body":
		return bufferedBodyHandler(w, r, splitURL[4])
	default:
		return xssTemplate(w, r, pd)
	}
}
