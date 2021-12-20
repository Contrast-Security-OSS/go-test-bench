package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const DefaultAddr = "localhost:8080"

func main() {
	log.Println("Starting...")
	addr := flag.String("addr", DefaultAddr, "`host:port` to access the listening server")
	flag.Parse()

	client := http.DefaultClient

	reqs, err := unsafeRequests(*addr)
	if err != nil {
		log.Fatalf("failed to generate requests: %s", err)
	}

	for _, r := range reqs {
		log.Printf("sending request: %s", r.URL.String())
		res, err := client.Do(r)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != 200 {
			log.Fatalf("unsuccessful response: %d", res.StatusCode)
		}
		log.Printf("route exercised: %s", r.URL.String())
	}

	log.Println("All routes successfully exercised")
}

// unsafeRequests generates a list of requests for all unsafe
// endpoints in the gin app.
func unsafeRequests(addr string) ([]*http.Request, error) {
	var reqs []*http.Request
	type genFunc func(addr string) ([]*http.Request, error)
	for _, gen := range []genFunc{
		unsafeCmdInjections,
		unsafePathTraversals,
		unsafeSQLInjections,
		unsafeSSRF,
		unsafeUnvalidatedRedirect,
		unsafeXSS,
	} {
		newReqs, err := gen(addr)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, newReqs...)
	}
	return reqs, nil
}

// addUserInput adds user controllable data to the request r.
// The data type can be configured with inputType. If inputType is not
// supported, the program exits.
// You can also specify the key and value of the data to be added to
// the request. The key "input" and value "fake-user-input"
// are used by default.
func addUserInput(inputType string, r *http.Request, key, value string) {
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
		r.URL.Path = path.Join(r.URL.Path, url.PathEscape(value))
	case "response":
		// BUG: This endpoint doesn't actually read a response.
		// For now, just add a header since it's quick
		r.Header.Set(key, value)
	default:
		log.Fatalf("unknown input type: %s", inputType)
	}
}
