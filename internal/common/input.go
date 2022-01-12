package common

import (
	"net/http"
	"strings"
)

//INPUT used as the standard input name across the project and forms
const INPUT = "input"

//GetUserInput returns the first value found in the request with the given key
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
func GetUserInput(r *http.Request) string {
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

	// GetPathValue is not included because we need to have the positions already defined
	// Currently it is used directly in XSS parsing <script> ... </script> from the query path
	// TODO - need to update the logic for the positional parameters

	return ""
}

//GetParamValue returns the input value for the given key of a GET request query
func GetParamValue(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

//GetPathValue returns the input value as a string included in the URL - e.g. /<script>.....</script>/....
//
// since the path is split by "/" there is a need to combine multiple pieces into one in order to get the full
// value accordingly, positions - holds the indices of the split string to concatenate
func GetPathValue(r *http.Request, positions ...int) string {
	splitURL := strings.Split(r.URL.Path, "/")
	var param []string
	for _, v := range positions {
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
	return r.Header.Get(key)
}
