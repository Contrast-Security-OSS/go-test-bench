package common

import (
	"fmt"
	"html/template"
	"log"
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
func writeCurlCmds(addr, base, url, input, payload string, extraargs ...string) template.HTML {
	var out []string
	modes := []struct{ name, frag string }{
		{"unsafe", "unsafe"},
		{"safe", "safe"},
		{"no-op", "noop"},
	}

	if len(payload) > 0 {
		switch {
		case input == "headers":
			extraargs = append(extraargs, `-H`, `"input: `+payload+`"`)
		case input == "headers-json":
			extraargs = append(extraargs,
				"-X", "POST",
				`-H`, `"Content-Type: application/json"`,
				"\\\n   ", //line break for better display
				`-H`, `"credentials:{\"username\":\"`+payload+`\",\"password\":\"12345Pass\"}"`)
		case input == "cookies":
			extraargs = append(extraargs,
				"-X", "POST",
				"-b", `"input=`+payload+`"`,
			)
		default:
			log.Fatalf("writeCurlCmds: unknown input type %q", input)
		}
	}
	args := strings.Join(extraargs, " ")
	for _, m := range modes {
		fullUrl := fmt.Sprintf("http://%s%s/%s/%s/%s", addr, base, url, input, m.frag)
		out = append(out, fmt.Sprintf("<p>%s</p><pre class=%q>curl %s \\\n    %s</pre>", m.name, m.frag, fullUrl, args))
	}
	return template.HTML(strings.Join(out, "\n"))
}

// do we use curl for this input type?
//
// current input types ( * --> needs curl ):
//
//       body
//     * cookies
//     * headers
//     * headers-json
//       params
//       query
//
func needsCurl(input string) bool {
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