// Package servetest contains test logic shared between multiple frameworks.
package servetest

import (
	"sort"
	"strings"
	"testing"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/google/go-cmp/cmp"
)

type RouteInputs map[string][]string

// DefaultRouteInputs defines the inputs used by each route. Exported so
// that it can be modified if a framework deviates from the default.
func DefaultRouteInputs() RouteInputs {
	return map[string][]string{
		"cmdInjection":        {"query", "cookies"},
		"pathTraversal":       {"query", "buffered-query", "headers", "body"},
		"xss":                 {"query", "buffered-query", "params", "body", "buffered-body", "response"},
		"sqlInjection":        {"headers-json", "query", "body"},
		"ssrf":                {"query", "path"},
		"unvalidatedRedirect": {"query"},
	}
}

// TestRouteData is called from a framework-specific test to check that
// input types are correct. rmap and/or routeInputs may be nil, in which
// case the default values are used.
func TestRouteData(t *testing.T, rmap common.RouteMap, routeInputs RouteInputs) {
	var routeNames []string
	if rmap == nil {
		rmap = common.GetRouteMap()
	}
	for k := range rmap {
		routeNames = append(routeNames, strings.TrimLeft(k, "/"))
	}
	sort.Strings(routeNames)
	for _, n := range routeNames {
		r, ok := rmap[n]
		if ok && (len(r.Sinks) == 0 || len(r.Sinks[0].Name) == 0) {
			delete(rmap, n)
		}
	}

	if routeInputs == nil {
		routeInputs = DefaultRouteInputs()
	}
	if len(routeInputs) != len(rmap) {
		var testNames []string
		routeNames = nil
		for k := range rmap {
			routeNames = append(routeNames, strings.TrimLeft(k, "/"))
		}
		sort.Strings(routeNames)
		for k := range routeInputs {
			testNames = append(testNames, k)
		}
		sort.Strings(testNames)
		t.Errorf("mismatch between number of routes and test data:\nwant %v\n got %v", testNames, routeNames)
	}
	for k, v := range rmap {
		k = strings.TrimLeft(k, "/")
		t.Run(k, func(t *testing.T) {
			types, ok := routeInputs[k]
			if !ok {
				t.Fatalf("missing test data for %s", k)
			}
			sort.Strings(v.Inputs)
			sort.Strings(types)
			if diff := cmp.Diff(v.Inputs, types); len(diff) > 0 {
				t.Errorf("input type mismatch for %s:\n%s", k, diff)
			}
		})
	}
}
