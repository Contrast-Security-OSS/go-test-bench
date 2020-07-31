package xxe

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/utils"
)

func encodingXMLHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, mode string) (template.HTML, bool) {
	var command *exec.Cmd
	var userInput string

	switch mode {
	case "safe":
		command = exec.Command("echo", userInput)

	case "unsafe":
		command = exec.Command("sh", "-c", "echo "+userInput)

	case "noop":
		return template.HTML("NOOP"), false

	default:
		log.Fatal("Error running encodingXMLHandler. No option passed")
	}

	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}

	return template.HTML(out.String()), false
}

func xxeTemplate(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	var buf bytes.Buffer
	templates := template.Must(template.ParseFiles("./views/partials/safeButtons.gohtml", "./views/pages/xxe.gohtml", "./views/partials/ruleInfo.gohtml"))
	err := templates.ExecuteTemplate(&buf, "xxe", routeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true

}

// Handler is the API handler for XXE
func Handler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	fmt.Println(splitURL)
	switch splitURL[2] {

	case "default":
		return encodingXMLHandler(w, r, routeInfo, splitURL[len(splitURL)-1])
	case "":
		return xxeTemplate(w, r, routeInfo)
	default:
		log.Fatal("XXE Handler reached incorrectly")
		return "", false
	}
}
