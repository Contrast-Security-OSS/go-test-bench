package xss

import (
	"bytes"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"log"
	"net/url"
)

// RegisterRoutes is to be called to add the routes in this package to common.AllRoutes.
func RegisterRoutes(frameworkSinks ...*common.Sink) {
	sinks := frameworkSinks
	common.Register(common.Route{
		Name:     "Reflected XSS",
		Link:     "https://www.owasp.org/index.php/Cross-site_Scripting_(XSS)#Stored_and_Reflected_XSS_Attacks",
		Base:     "xss",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "buffered-query", "params", "body", "buffered-body"},
		Sinks:    sinks,
		Payload:  "<script>alert(1);</script>",
	})
}

func CommonHandler(safety common.Safety, payload string) (data string) {
	switch safety {
	case "safe":
		payload = url.QueryEscape(payload)
	case "unsafe":
	case "noop":
		payload = "NOOP"
	default:
		log.Fatal("Error running CommonHandler. No option passed")
	}
	//execute input script
	return payload

}

// CommonBufferedHandler used as a handler which uses bytes.Buffer for source input ignoring the user input
func CommonBufferedHandler(safety common.Safety, payload string) (data string) {
	var buf bytes.Buffer
	buf.WriteString(payload)

	switch safety {
	case "safe":
		payload = string(buf.Bytes())
		payload = url.QueryEscape(payload)
	case "unsafe":
		payload = string(buf.Bytes())
	case "noop":
		payload = "NOOP"
	default:
		log.Fatal("Error running CommonBufferedHandler. No option passed")
	}

	return payload
}
