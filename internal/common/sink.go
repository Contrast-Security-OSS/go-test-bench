package common

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// A SanitizerFn sanitizes the input string
type SanitizerFn func(string) string

// Sink is a struct that identifies the name of the sink, the associated URL and
// the HTTP method
type Sink struct {
	Name string
	URL  string

	// if nil, a generic handler is used and VulnerableFnWrapper and Sanitizer must
	// both be set
	Handler HandlerFn

	// a function that renders input safe; only used by the generic handler and only
	// when 'safe' mode is requested.
	//
	// for example: url.QueryEscape
	Sanitize SanitizerFn

	// the vulnerable function which may recieve unsanitized input. Handler must be
	// nil when this is set.
	VulnerableFnWrapper VulnerableFnWrapper
}

func (s *Sink) String() string {
	if len(s.Name) == 0 || s.Name == "_" {
		return ""
	}
	return fmt.Sprintf("%s: %s", s.Name, path.Join("...", s.URL))
}

// AddPayloadToRequest adds user controllable data to the request r.
// The data type can be configured with inputType. If inputType is not
// supported, the program exits.
// You can also specify the key and value of the data to be added to
// the request. The key "input" and value "fake-user-input"
// are used by default.
func (s *Sink) AddPayloadToRequest(req *http.Request, inputType, key, payload string) {
	if len(key) == 0 {
		key = "input"
	}
	if payload == "" {
		payload = "fake-user-input"
	}
	switch inputType {
	case "query", "buffered-query":
		q := req.URL.Query()
		q.Add(key, payload)
		req.URL.RawQuery = q.Encode()
	case "body", "buffered-body":
		form := make(url.Values)
		form.Set(key, payload)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	case "cookies":
		req.AddCookie(&http.Cookie{
			Name:  key,
			Value: payload,
		})
	case "headers":
		req.Header.Set(key, payload)
	case "params":
		req.URL.Path = path.Join(req.URL.Path, payload)
	case "response":
		// BUG: This endpoint doesn't actually read a response.
		// For now, just add a header since it's quick
		req.Header.Set(key, payload)
	case "headers-json":
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		creds.Username = payload
		data, err := json.Marshal(creds)
		if err != nil {
			log.Fatalf("failed to marshal JSON object: %s", err)
		}
		// special case for sqli (the only user of headers-json): use 'credentials'
		// for the key, rather than whatever is provided
		req.Header.Set("credentials", string(data))
	default:
		log.Fatalf("unknown input type: %s", inputType)
	}
}