package commontest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// AddUserInput adds user controllable data to the request r.
// The data type can be configured with inputType. If inputType is not
// supported, the program exits.
// You can also specify the key and value of the data to be added to
// the request. The key "input" and value "fake-user-input"
// are used by default.
func AddUserInput(inputType string, r *http.Request, key, value string) {
	if key == "" {
		key = "input"
	}
	if value == "" {
		value = "fake-user-input"
	}
	switch inputType {
	case "query", "buffered-query":
		q := r.URL.Query()
		q.Add(key, value)
		r.URL.RawQuery = q.Encode()
	case "body", "buffered-body":
		v := make(url.Values)
		v.Set(key, value)
		r.Body = io.NopCloser(strings.NewReader(v.Encode()))
	case "cookies":
		r.AddCookie(&http.Cookie{
			Name:  key,
			Value: value,
		})
	case "headers":
		r.Header.Set(key, value)
	case "params":
		r.URL.Path = path.Join(r.URL.Path, value)
	case "response":
		// BUG: This endpoint doesn't actually read a response.
		// For now, just add a header since it's quick
		r.Header.Set(key, value)
	default:
		log.Fatalf("unknown input type: %s", inputType)
	}
}

func addHeadersJSONCreds(r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	creds.Username = "Robert'; DROP TABLE Students;--"
	data, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("failed to marshal JSON object: %s", err)
	}
	r.Header.Set("credentials", string(data))
}
