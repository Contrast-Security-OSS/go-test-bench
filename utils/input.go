package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//INPUT used as the standard input name across the project and forms
const INPUT = "input"

//GetQueryInput returns the input value for the given key of a GET request query
func GetQueryInput(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

//GetParameterInput returns the input value as a string included in the URL - e.g. /<script>.....</script>/....
func GetParameterInput(r *http.Request, positions ...int) string {
	splitURL := strings.Split(r.URL.Path, "/")
	var param []string
	for _, v := range positions {
		param = append(param, splitURL[v])
	}
	return strings.Join(param, "/")
}

//PostInput returns the input value from the POST body of the request
func PostInput(r *http.Request, key string) (string, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	//the request body should only have the user input
	inputs, _ := url.QueryUnescape(string(b)) //Couldn't find whether POST request inputs get sanitizied ALWAYS.
	splitInput := strings.Split(inputs, key+"=")

	return splitInput[1], nil
}

//FormValueInput returns the input value for the given key from the submitted form
func FormValueInput(r *http.Request, key string) string {
	return r.FormValue(key)
}

//CookieInput returns the input value for the given cookie
func CookieInput(r *http.Request, key string) string {
	cookie, err := r.Cookie(key)
	if err != nil {
		log.Printf("Could not get %s cooke value, error:%s\n", key, err.Error())
	}

	return cookie.Value
}

//HeaderInput returns the input value from the given header
func HeaderInput(r *http.Request, key string) string {
	return r.Header.Get(key)
}
