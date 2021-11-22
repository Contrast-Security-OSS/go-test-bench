package pathtraversal

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

func pathTTemplate(w http.ResponseWriter, r *http.Request, data common.Parameters) (template.HTML, bool) {
	return "pathTraversal.gohtml", true
}

func readFileHandler(w http.ResponseWriter, r *http.Request, method string, inputs string, isBuffered bool) (template.HTML, bool) {
	log.Printf("Executing pathTraversal readFile ; method = %s ; use of bytes.Buffer = %t\n", method, isBuffered)
	var data string
	var err error
	switch method {
	case "safe":
		inputs = url.QueryEscape(inputs)
		if isBuffered {
			data, err = bufferedReadFile(inputs)
		} else {
			data, err = readFile(inputs)
		}

		if err != nil {
			return template.HTML("Congrats, you are safe! Error from readFile: " + err.Error()), false
			//fmt.Fprintf(w, "Congrats, you are safe! Error from readFile: ")
			//http.Error(w, err.Error(),500)
		} else if data == "" {
			return template.HTML("Congrats, you are safe! No data was read."), false
		}
	case "unsafe":
		if isBuffered {
			data, err = bufferedReadFile(inputs)
		} else {
			data, err = readFile(inputs)
		}
	default:
		data = "INVALID URL"
	}
	return template.HTML(data), false
}

// bufferedReadFile read the given file using bytes.Buffer
func bufferedReadFile(filename string) (string, error) {
	fr, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer fr.Close()

	var buf bytes.Buffer
	buf.ReadFrom(fr)

	return buf.String(), nil
}

// readFile read the given file using ioutil.ReadFile
func readFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	data := string(content)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if data == "" || err != nil {
		data = "Done!"
	} else {
		data = string(content)
	}
	return data, nil
}

func writeFileHandler(w http.ResponseWriter, r *http.Request, method string, inputs string, isBuffered bool) (template.HTML, bool) {
	log.Printf("Executing pathTraversal writeFile ; method =  %s ; use of bytes.Buffer = %t\n", method, isBuffered)
	var data, message string
	var err error
	switch method {
	case "safe":
		inputs = url.QueryEscape(inputs)
		if isBuffered {
			message = "buffered pathTraversal"
			err = bufferedWriteFile(inputs, message)
		} else {
			message = "pathTraversal"
			err = writeFile(inputs, message)
		}
		data = message
		if err != nil {
			return template.HTML("Congrats, you are safe! Error from writeFile: " + err.Error()), false
		} else if data == "" {
			return template.HTML("Congrats, you are safe! No data was written."), false
		} else { //something was written!
			data = "Wrote \"" + message + "\" to: " + inputs
			log.Printf("%s was written with pathTraversal\n", inputs)
		}
	case "unsafe":
		if isBuffered {
			message = "buffered pathTraversal"
			err = bufferedWriteFile(inputs, message)
		} else {
			message = "pathTraversal"
			err = writeFile(inputs, message)
		}
		data = message
		if data == "" || err != nil {
			data = "Done!"
		} else {
			data = "Wrote \"" + message + "\" to: " + inputs
		}
	default:
		data = "INVALID URL"
	}
	return template.HTML(data), false
}

// bufferedWriteFile write message content in the given file using bytes.Buffer
func bufferedWriteFile(filename, message string) error {
	var buf bytes.Buffer
	fmt.Fprint(&buf, message)

	fr, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fr.Close()

	_, err = buf.WriteTo(fr)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// bufferedWriteFile write message content in the given file using ioutil.WriteFile
func writeFile(filename, message string) error {
	err := ioutil.WriteFile(filename, []byte(message), 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Handler is the API handler for path traversal
func Handler(w http.ResponseWriter, r *http.Request, pd common.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	if len(splitURL) < 4 {
		return pathTTemplate(w, r, pd)
	}
	if splitURL[2] != "body" && splitURL[2] != "headers" && splitURL[2] != "query" && splitURL[2] != "buffered-query" {
		return template.HTML("INVALID URL"), false
	}
	if splitURL[4] == "noop" {
		return template.HTML("NOOP"), false
	}

	isBuffered := strings.Contains(splitURL[2], "buffered")

	userInput := common.GetUserInput(r)

	switch splitURL[3] {
	case "ioutil.ReadFile":
		return readFileHandler(w, r, splitURL[4], userInput, isBuffered)
	case "ioutil.WriteFile":
		return writeFileHandler(w, r, splitURL[4], userInput, isBuffered)
	default: //should be an error instead
		return template.HTML("INVALID URL"), false
	}
}
