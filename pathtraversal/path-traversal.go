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
			log.Fatal("ERROR: pathTraversal readFile safe button not properly implemented")
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
			//if you want to see the actual data, remove this else statement
			data = "stuff"
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
		message := []byte("")
		err := ioutil.WriteFile(inputs, message, 0644)
		if err != nil {
			return template.HTML("Congrats, you are safe! Error from writeFile: " + err.Error()), false
		} else if string(message) == "" {
			return template.HTML("Congrats, you are safe! No data was written."), false
		} else { //something was written!
			log.Fatal("ERROR: pathTraversal writeFile safe button not properly implemented")
		}
	case "unsafe":
		message := []byte("")
		err := ioutil.WriteFile(inputs, message, 0644)
		data = string(message)
		if err != nil {
			log.Println(err)
		}
		if data == "" || err != nil {
			data = "Done!"
		} else {
			data = "Wrote to: " + data
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
	if splitURL[4] == "noop" {
		return template.HTML("NOOP"), false
	}
	var inputs string
	switch splitURL[2] {
	case "body":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return template.HTML(err.Error()), false
		}
		//the request body should only have the user input
		inputs, err = url.QueryUnescape(string(b)) //Couldn't find whether POST request inputs get sanitizied ALWAYS.
		splitInput := strings.Split(inputs, "input=")
		inputs = splitInput[1]
		if err != nil {
			return template.HTML(err.Error()), false
		}
	case "headers":
		inputs = r.Header["Input"][0]
	case "query":
		inputs = r.URL.Query().Get("input")
	default:
		return template.HTML("INVALID URL"), false
	}
	switch splitURL[3] {
	case "ioutil.ReadFile":
		return readFileHandler(w, r, splitURL[4], inputs)
	case "ioutil.WriteFile":
		return writeFileHandler(w, r, splitURL[4], inputs)
	default: //should be an error instead
		return template.HTML("INVALID URL"), false
	}
}
