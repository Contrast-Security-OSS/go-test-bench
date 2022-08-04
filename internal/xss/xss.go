package xss

import (
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"net/url"
)

// RegisterRoutes is to be called to add the routes in this package to common.AllRoutes.
func RegisterRoutes(frameworkSinks ...*common.Sink) {
	sinks := []*common.Sink{
		{
			Name:     "reflectedXss",
			Sanitize: url.PathEscape,
			VulnerableFnWrapper: func(_ interface{}, payload string) (data string, raw bool, err error) {
				return payload, true, nil
			},
			RawMime: "text/html",
		},
	}

	if len(frameworkSinks) > 0 {
		sinks = append(sinks, frameworkSinks...)
	}

	common.Register(common.Route{
		Name:     "Reflected XSS",
		Link:     "https://www.owasp.org/index.php/Cross-site_Scripting_(XSS)#Stored_and_Reflected_XSS_Attacks",
		Base:     "xss",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "buffered-query", "params", "body", "buffered-body", "response"},
		Sinks:    sinks,
		Payload:  "<html><img src=a onerror=alert(1)>",
	})
}
