package common

import (
	"fmt"
	"html/template"
	"strings"
)

// FuncMap returns a map of functions for use from templates.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		// lowercase the input
		"tolower": strings.ToLower,
		// transforms an "input" to a title suitable for display. For example, "headers-json" becomes "Headers - JSON"
		"inputTitle": inputTitle,
		// writes html with unsafe/safe/noop curl commands
		"curl": writeCurlCmds,
		// determine if input type is for curl
		"needsCurl": needsCurl,
	}
}

// use in template to generate unsafe/safe/no-op curl commands
// call once to generate 3 commands
// {{ curl $addr $base .URL "headers-json" "POST" "-H \"Content-Type: application/json\" \\\n    -H \"credentials:{...}\""}}
// -> curl http://localhost:8080/sqlInjection/sqlite3.exec/headers-json/unsafe -X POST -H ...
func writeCurlCmds(addr, base, url, input, method, args string) template.HTML {
	var out []string
	modes := []struct{ name, frag string }{
		{"unsafe", "unsafe"},
		{"safe", "safe"},
		{"no-op", "noop"},
	}
	for _, m := range modes {
		out = append(out, fmt.Sprintf("<p>%s</p><pre>curl http://%s%s/%s/%s/%s \\\n    -X %s %s</pre>", m.name, addr, base, url, input, m.frag, method, args))
	}
	return template.HTML(strings.Join(out, "\n"))
}

// do we use curl for this input type?
func needsCurl(input string) bool {
	input = strings.ToLower(input)
	switch {
	case strings.Contains(input, "header"):
		return true
	case strings.Contains(input, "cookie"):
		return true
	default:
		return false
	}
}

// transforms an "input" to a title suitable for display. For example, "headers-json" becomes "Headers - JSON"
func inputTitle(input string) string {
	if strings.Contains(input, " ") {
		panic(fmt.Sprintf(".Input cannot contain spaces: %q", input))
	}
	input = strings.ToLower(input)
	title := strings.ReplaceAll(input, "-", " - ")
	title = strings.ReplaceAll(title, "json", "JSON")
	title = strings.Title(title)
	return title
}
