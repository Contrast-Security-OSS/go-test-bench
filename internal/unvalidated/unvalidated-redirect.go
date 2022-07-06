package unvalidated

import (
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

func RegisterRoutes(frameworkSinks ...*common.Sink) {
	sinks := frameworkSinks
	common.Register(common.Route{
		Name:     "Unvalidated Redirect",
		Link:     "https://cheatsheetseries.owasp.org/cheatsheets/Unvalidated_Redirects_and_Forwards_Cheat_Sheet.html",
		Base:     "unvalidatedRedirect",
		Products: []string{"Assess"},
		Inputs:   []string{"query"},
		Sinks:    sinks,
		Payload:  "http://example.com",
	})
}
