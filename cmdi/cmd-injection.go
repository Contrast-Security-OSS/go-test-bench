package cmdi

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/utils"
)

var templates = template.Must(template.ParseFiles(
	"./views/partials/safeButtons.gohtml",
	"./views/pages/commandInjection.gohtml",
	"./views/partials/ruleInfo.gohtml",
))

func osExecHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, splitURL []string) (template.HTML, bool) {
	var command *exec.Cmd
	mode := splitURL[len(splitURL)-1]

	userInput := utils.GetUserInput(r)

	switch mode {
	case "safe":
		command = exec.Command("echo", userInput)

	case "unsafe":
		command = exec.Command(userInput)

	case "noop":
		return template.HTML("NOOP"), false

	default:
		log.Fatal("Error running execHandler. No option passed")
	}

	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		log.Printf("Could not run command in %s: err = %s", mode, err)
	}

	return template.HTML(out.String()), false
}

// Feed the commandInjection template with page parameters
func cmdiTemplate(w http.ResponseWriter, r *http.Request, params utils.Parameters) (template.HTML, bool) {
	var buf bytes.Buffer

	err := templates.ExecuteTemplate(&buf, "commandInjection", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true

}

// Handler is the API handler for command injection
func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	switch splitURL[2] {

	case "osExec":
		return osExecHandler(w, r, pd.Rulebar[pd.Name], splitURL)
	case "":
		return cmdiTemplate(w, r, pd)

	default:
		log.Fatal("commandInjection Handler reached incorrectly")
		return "", false
	}
}
