package ssrf

//if we want to get this package we just call "ssrf" as an import now
import (
	//"fmt"
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/utils"
)

var templates = template.Must(template.ParseFiles(
	"./views/pages/ssrf.gohtml",
	"./views/partials/ruleInfo.gohtml",
))

// Handler is the API handler for SSRF
func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	routeInfo := pd.Rulebar[pd.Name]
	if r.URL.Path == "/ssrf/" { //or "/ssrf"
		return bodyHandler(w, r, routeInfo)
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
	switch method {
	case "query":
		inputs := r.URL.Query().Get("input")
		if inputs != "" {
			res, err = http.Get("http://example.com?input=" + inputs)
		} else {
			res, err = http.Get("http://example.com")
		}
	case "path":
		url := r.URL.Query().Get("input")
		if url == "" {
			url = "example.com"
		}
		res, err = http.Get("http://" + url)
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

func bodyHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "ssrf", &routeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true
}
