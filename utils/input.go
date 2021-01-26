package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//INPUT used as the standard input name across the project and forms
const INPUT = "input"

//GetParamValue returns the input value for the given key of a GET request query
func GetParamValue(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

//GetPathValue returns the input value as a string included in the URL - e.g. /<script>.....</script>/....
func GetPathValue(r *http.Request, positions ...int) string {
	splitURL := strings.Split(r.URL.Path, "/")
	var param []string
	for _, v := range positions {
		param = append(param, splitURL[v])
	}
	return strings.Join(param, "/")
}

//GetPostBody returns the input value from the POST body of the request
func GetPostBody(r *http.Request, key string) (string, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//GetFormValue returns the input value for the given key from the submitted form
func GetFormValue(r *http.Request, key string) string {
	return r.FormValue(key)
}

//GetCookieValue returns the input value for the given cookie
func GetCookieValue(r *http.Request, key string) string {
	cookie, err := r.Cookie(key)
	if err != nil {
		log.Printf("Could not get %s cooke value, error:%s\n", key, err.Error())
	}

	return cookie.Value
}

//GetHeaderValue returns the input value from the given header
func GetHeaderValue(r *http.Request, key string) string {
	return r.Header.Get(key)
}
