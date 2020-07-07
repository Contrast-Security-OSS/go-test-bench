package xxe

import (
	"net/http"
	"html/template"
	"bytes"
	"strings"
	"log"
	"os/exec"
	"fmt"
	//"encoding/xml"
	utils "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/utils"
)

var templates = template.Must(template.ParseFiles("./views/partials/safeButtons.gohtml","./views/pages/xxe.gohtml", "./views/partials/ruleInfo.gohtml"))


func encodingXMLHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, mode string) (template.HTML, bool) {	
	var command *exec.Cmd
	var user_input string

	switch mode {
		case "safe":
			command = exec.Command("echo", user_input)

		case "unsafe":
			command = exec.Command("sh", "-c", "echo " + user_input)

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

func defaultHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "xxe", routeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true 

}


func Handler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	split_url := strings.Split(r.URL.Path, "/")
	fmt.Println(split_url)
	switch split_url[2] {
		
		case "default":
			return encodingXMLHandler(w,r,routeInfo,split_url[len(split_url) - 1])
		case "":
			return defaultHandler(w, r, routeInfo)
		default:
			log.Fatal("XXE Handler reached incorrectly")
			return "", false
		}
}