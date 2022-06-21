package ssrf

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// RegisterRoutes is to be called to add the routes in this package to common.AllRoutes.
func RegisterRoutes(frameworkSinks ...*common.Sink) {
	sinks := []*common.Sink{
		{
			Name:    "net/http",
			URL:     "http",
			Handler: httpHandler,
		},
	}
	sinks = append(sinks, frameworkSinks...)
	common.Register(common.Route{
		Name:     "Server Side Request Forgery",
		Link:     "https://owasp.org/www-community/attacks/Server_Side_Request_Forgery",
		Base:     "ssrf",
		Products: []string{"Assess"},
		Inputs:   []string{"query", "params"},
		Sinks:    sinks,
		Payload:  "http://example.com",
	})
}

func httpHandler(safety common.Safety, payload string, opaque interface{}) (data, mime string, status int) {
	mime = "text/plain"
	if len(payload) == 0 {
		data = "payload required but not provided"
		log.Println(data)
		return data, mime, http.StatusBadRequest
	}
	if u, err := url.Parse(payload); err != nil {
		data = fmt.Sprintf("can't parse url %q: %s", payload, err)
		log.Println(data)
		return data, mime, http.StatusBadRequest
	} else if len(u.Scheme) == 0 {
		log.Printf("missing scheme in url %q, adding http:// prefix...", payload)
		payload = "http://" + payload
	}
	switch safety {
	case common.Unsafe:
		resp, err := http.Get(payload)
		if err != nil {
			data = fmt.Sprintf("failed to GET %q: %s", payload, err)
			log.Println(data)
			return data, mime, http.StatusInternalServerError
		}
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			data = fmt.Sprintf("failed to read body from %q (upstream response %s): %s", payload, resp.Status, err)
			log.Println(data)
			return data, mime, http.StatusInternalServerError
		}
		if ct := resp.Header.Get("Content-Type"); len(ct) > 0 {
			mime = ct
		}
		return string(bdy), mime, resp.StatusCode
	case common.NOOP:
		return fmt.Sprintf("%s %s", safety, payload), mime, http.StatusOK
	default:
		return fmt.Sprintf("unhandled safety type %q", safety), mime, http.StatusInternalServerError
	}
}
