package pathtraversal

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/utils"
)

var templates = template.Must(template.ParseFiles(
	"./views/partials/safeButtons.gohtml",
	"./views/pages/pathTraversal.gohtml",
	"./views/partials/ruleInfo.gohtml",
))

func pathTTemplate(w http.ResponseWriter, r *http.Request, data utils.Parameters) (template.HTML, bool) {
	var buf bytes.Buffer

	err := templates.ExecuteTemplate(&buf, "pathTraversal", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true
}

func readFileHandler(w http.ResponseWriter, r *http.Request, method string, inputs string) (template.HTML, bool) {
	log.Println("Executing pathTraversal readFile ; method = " + method)
	var data string
	switch method {
	case "safe":
		inputs = url.QueryEscape(inputs)
		content, err := ioutil.ReadFile(inputs)
		if err != nil {
			return template.HTML("Congrats, you are safe! Error from readFile: " + err.Error()), false
			//fmt.Fprintf(w, "Congrats, you are safe! Error from readFile: ")
			//http.Error(w, err.Error(),500)
		} else if string(content) == "" {
			return template.HTML("Congrats, you are safe! No data was read."), false
		} else {
			data = string(content)
		}
	case "unsafe":
		content, err := ioutil.ReadFile(inputs)
		data = string(content)
		if err != nil {
			log.Println(err)
		}
		if data == "" || err != nil {
			data = "Done!"
		} else {
			data = string(content)
		}
	default:
		data = "INVALID URL"
	}
	return template.HTML(data), false
}

func writeFileHandler(w http.ResponseWriter, r *http.Request, method string, inputs string) (template.HTML, bool) {
	log.Println("Executing pathTraversal writeFile ; method = " + method)
	var data string
	switch method {
	case "safe":
		inputs = url.QueryEscape(inputs)
		message := []byte("pathTraversal")
		err := ioutil.WriteFile(inputs, message, 0644)
		if err != nil {
			return template.HTML("Congrats, you are safe! Error from writeFile: " + err.Error()), false
		} else if string(message) == "" {
			return template.HTML("Congrats, you are safe! No data was written."), false
		} else { //something was written!
			data = "Wrote \"" + string(message) + "\" to: " + inputs
			log.Printf("%s was written with pathTraversal\n", inputs)
		}
	case "unsafe":
		message := []byte("pathTraversal")
		err := ioutil.WriteFile(inputs, message, 0644)
		data = string(message)
		if err != nil {
			log.Println(err)
		}
		if data == "" || err != nil {
			data = "Done!"
		} else {
			data = "Wrote \"" + data + "\" to: " + inputs
		}
	default:
		data = "INVALID URL"
	}
	return template.HTML(data), false
}

//Handler is the API handler for path traversal
func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	if len(splitURL) < 4 {
		return pathTTemplate(w, r, pd)
	}
	if splitURL[2] != "body" && splitURL[2] != "headers" && splitURL[2] != "query" {
		return template.HTML("INVALID URL"), false
	}
	if splitURL[4] == "noop" {
		return template.HTML("NOOP"), false
	}

	userInput := utils.GetUserInput(r)

	switch splitURL[3] {
	case "ioutil.ReadFile":
		return readFileHandler(w, r, splitURL[4], userInput)
	case "ioutil.WriteFile":
		return writeFileHandler(w, r, splitURL[4], userInput)
	default: //should be an error instead
		return template.HTML("INVALID URL"), false
	}
}
