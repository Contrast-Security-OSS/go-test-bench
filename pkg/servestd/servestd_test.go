package servestd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/google/go-cmp/cmp"
)

func TestParseTemplates(t *testing.T) {
	err := common.ParseViewTemplates()
	if err != nil {
		t.Error(err)
	}
}

//ensure template is populated
func TestServedTemplate(t *testing.T) {
	Setup()
	t.Cleanup(common.Reset)

	srvr := httptest.NewServer(nil)
	defer srvr.Close()
	Pd.Addr = srvr.URL
	for _, route := range Pd.Rulebar {
		if len(route.Sinks) == 0 || len(route.Sinks[0].Name) == 0 {
			continue
		}
		t.Run(route.Name, func(t *testing.T) {
			url := srvr.URL + route.Base + "/"
			t.Logf("url %s", url)

			resp, err := http.Get(url)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			const findOut = "To find out more about "
			if !strings.Contains(string(body), findOut+route.Name) {
				lines := strings.Split(string(body), "\n")
				var i int
				var found bool
				for i = range lines {
					if strings.Contains(lines[i], findOut) {
						found = true
						break
					}
				}
				if found {
					t.Errorf("missing command name at %d:\n--  %s\n==> %s\n-- %s\n", i, lines[i-1], lines[i], lines[i+2])
				} else {
					t.Errorf("missing command name in\n%s", body)
				}
				var rb []string
				for _, v := range Pd.Rulebar {
					rb = append(rb, v.String())
				}
				t.Logf("\n%s", strings.Join(rb, "\n"))
			}
		})
	}
}

// changes here should be mirrored in the test in servegin
func TestRouteData(t *testing.T) {
	Setup()
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
