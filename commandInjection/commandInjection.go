package commandInjection

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"

	utils "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/utils"
)

var templates = template.Must(template.ParseFiles("./views/partials/safeButtons.gohtml", "./views/pages/commandInjection.gohtml", "./views/partials/ruleInfo.gohtml"))

func osExecHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, split_url []string) (template.HTML, bool) {
	var command *exec.Cmd
	var user_input string
	mode := split_url[len(split_url)-1]

	if user_input = r.FormValue("input"); user_input == "" {
		cookie, _ := r.Cookie("input")
		user_input = cookie.Value
	}


	switch mode {
		case "safe":
			command = exec.Command("echo", user_input)

		case "unsafe":
			command = exec.Command("sh", "-c", "echo " + user_input)

		case "noop":
			return template.HTML("NOOP"), false

		default:
			log.Fatal("Error running execHandler. No option passed")
	}

	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}

	return template.HTML(out.String()), false
}

func defaultHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "commandInjection", routeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true 

}


func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	split_url := strings.Split(r.URL.Path, "/")
	switch split_url[2] {
		
		case "osExec":
			return osExecHandler(w,r,pd.Rulebar[pd.Name], split_url)
		case "":
			return defaultHandler(w, r, pd.Rulebar[pd.Name])
		
		default:
			log.Fatal("commandInjection Handler reached incorrectly")
			return "", false
		}
}
