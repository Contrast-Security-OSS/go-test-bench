// Package shared contains functionality used by `cmd/swagger` at runtime and by
// `cmd/swagger/regen`. Previously, `regen` was generating code in a package it
// depended on - so any deletion or corruption of generated files in that package
// would cause future regen runs to fail. Moving this code to a different package
// avoids that issue.
package shared

import (
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/sqli"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/pathtraversal"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/ssrf"
)

// FilterInputTypes removes unsupported input types so they are not rendered.
// NOTE: this filtering only has an effect on routes using the generic template
// TODO(XXX): support things other than query and buffered-query
func FilterInputTypes(rmap common.RouteMap) {
	allowContains := []string{
		"query", // query, buffered-query

		// wildcard path parameters are not allowed in openapi v2 or v3:
		// https://github.com/OAI/OpenAPI-Specification/issues/892#issuecomment-281170254
		//
		// workaround is to set up as a single urlencoded param (see comments on above) -
		// but how likely is this in a real-world swagger app?
	}
	for path, rt := range rmap {
		for i := 0; i < len(rt.Inputs); {
			allow := false
			for _, a := range allowContains {
				if strings.Contains(rt.Inputs[i], a) {
					allow = true
					break
				}
			}
			if !allow {
				rt.Inputs = append(rt.Inputs[:i], rt.Inputs[i+1:]...)
				continue
			}
			i++
		}
		rmap[path] = rt
	}
}

// RegisterNewRoutes is called from Setup() and from the code regenerator
func RegisterNewRoutes() {
	cmdi.RegisterRoutes()
	sqli.RegisterRoutes()
	pathtraversal.RegisterRoutes()
	ssrf.RegisterRoutes()
}
