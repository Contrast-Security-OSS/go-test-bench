package servegin

import (
	"sort"
	"strings"
	"testing"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/google/go-cmp/cmp"
)

// changes here should be mirrored in the test in servestd
func TestRouteData(t *testing.T) {
	Setup("don't care")
	t.Cleanup(common.Reset)
	rmap := common.GetRouteMap()
	var routeNames []string
	for k := range rmap {
		routeNames = append(routeNames, strings.TrimLeft(k, "/"))
		sort.Strings(routeNames)
	}
	for _, n := range routeNames {
		r := rmap[n]
		if len(r.Sinks) == 0 || len(r.Sinks[0].Name) == 0 {
			delete(rmap, n)
		}
	}

	tests := map[string][]string{
		"cmdInjection":        {"query", "cookies"},
		"pathTraversal":       {"query", "buffered-query", "headers", "body"},
		"xss":                 {"query", "buffered-query", "params", "body", "buffered-body", "response"},
		"sqlInjection":        {"headers-json", "query", "body"},
		"ssrf":                {"query", "path"},
		"unvalidatedRedirect": {"query"},
	}
	if len(tests) != len(rmap) {
		var testNames []string
		routeNames = nil
		for k := range rmap {
			routeNames = append(routeNames, strings.TrimLeft(k, "/"))
			sort.Strings(routeNames)
		}
		for k := range tests {
			testNames = append(testNames, k)
		}
		sort.Strings(testNames)
		t.Errorf("mismatch between number of routes and test data:\nwant %v\n got %v", testNames, routeNames)
	}
	for k, v := range rmap {
		k = strings.TrimLeft(k, "/")
		t.Run(k, func(t *testing.T) {
			types, ok := tests[k]
			if !ok {
				t.Fatalf("missing test data for %s", k)
			}
			sort.Strings(v.Inputs)
			sort.Strings(types)
			if diff := cmp.Diff(v.Inputs, types); len(diff) > 0 {
				// if !reflect.DeepEqual(v.Inputs, types) {
				t.Errorf("input type mismatch for %s:\n%s", k, diff)
			}
		})
	}
}
