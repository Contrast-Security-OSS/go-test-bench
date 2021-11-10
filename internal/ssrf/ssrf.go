package ssrf

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// Handler is the API handler for SSRF
func Handler(w http.ResponseWriter, r *http.Request, pd common.Parameters) (template.HTML, bool) {
	if r.URL.Path == "/ssrf/" { //or "/ssrf"
		return bodyHandler(w, r)
	}

	//the library being used should be stored in sep[2], path or query in sep[3]
	sep := strings.Split(r.URL.Path, "/")
	switch sep[2] {
	case "default":
		httpHandler(w, r, sep[3])
	case "http":
		httpHandler(w, r, sep[3])
	case "request":
		requestHandler(w, r, sep[3])
	default:
		log.Println("THIS IS NOT A LIBRARY")
	}

	return template.HTML(""), false
}

func httpHandler(w http.ResponseWriter, r *http.Request, method string) {
	var res *http.Response
	var err error
	userInput := common.GetUserInput(r)
	switch method {
	case "query":
		if userInput != "" {
			res, err = http.Get("http://example.com?input=" + userInput)
		} else {
			res, err = http.Get("http://example.com")
		}
	case "path":
		if userInput == "" {
			userInput = "example.com"
		}
		res, err = http.Get("http://" + userInput)
		//don't want http.Redirect(w, r, "http://" + url, http.StatusFound)
	default:
		log.Println("THIS IS NOT A METHOD")
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request, method string) {
	// var res request.SugaredResp
	// var err error
	// var m monacoClient
	// switch method {
	// case "query":
	// 	inputs := r.URL.Query().Get("input")
	// 	if inputs != "" {
	// 		m.Params = map[string]string{"input": inputs}
	// 	}
	// 	m.URL = "https://example.com"
	// case "path":
	// 	m.URL = r.URL.Query().Get("input")
	// 	if m.URL == ""{
	// 		m.URL = "example.com"
	// 	}
	// 	m.URL = "http://" + m.URL
	// default:
	// 	log.Println("THIS IS NOT A METHOD")
	// 	return
	// }
	// m.Method = "GET"
	// client := request.Client{
	// 	URL:    m.URL,
	// 	Method: m.Method,
	// 	Params: m.Params,
	// }
	// res, err = client.Do()
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(),500)
	// 	return
	// }
	// _, err = w.Write(res.Data)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(),500)
	// }
}

func bodyHandler(w http.ResponseWriter, r *http.Request) (template.HTML, bool) {
	return "ssrf.gohtml", true
}
