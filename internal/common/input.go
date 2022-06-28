package common

import (
	"net/http"
	"strings"
)

//INPUT used as the standard input name across the project and forms
const INPUT = "input"

// GetUserInput returns the first value found in the request with the key 'input'.
//
// If none are found, it tries for a header with key 'credentials', and finally
// the last element in the url.
//
// the order of precedence when getting the result is:
//
// - query parameter
//
// - form value
//
// - cookie value
//
// - header value
//
// - credentials header
//
func GetUserInput(r *http.Request) (val string) {
	if value := GetParamValue(r, INPUT); value != "" {
		return value
	}

	if value := GetFormValue(r, INPUT); value != "" {
		return value
	}

	if value := GetCookieValue(r, INPUT); value != "" {
		return value
	}

	if value := GetHeaderValue(r, INPUT); value != "" {
		return value
	}

	if value := GetHeaderValue(r, "credentials"); value != "" {
		return value
	}
	return GetPathValue(r, -1)
}

//GetParamValue returns the input value for the given key of a GET request query
func GetParamValue(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetPathValue returns element(s) from the given position(s) in the url,
// joined with '/'. Negative positions are allowed and start at the right.
func GetPathValue(r *http.Request, positions ...int) string {
	splitURL := strings.Split(r.URL.Path, "/")
	var param []string
	for _, v := range positions {
		if v < 0 {
			v = len(splitURL) + v
		}
		param = append(param, splitURL[v])
	}
	return strings.Join(param, "/")
}

//GetFormValue returns the input value for the given key from the submitted form
func GetFormValue(r *http.Request, key string) string {
	return r.FormValue(key)
}

//GetCookieValue returns the input value for the given cookie
func GetCookieValue(r *http.Request, key string) string {
	cookie, err := r.Cookie(key)
	if err != nil {
		return ""
	}

	return cookie.Value
}

//GetHeaderValue returns the input value from the given header
func GetHeaderValue(r *http.Request, key string) string {
	res := r.Header.Get(key)
	return res
}
