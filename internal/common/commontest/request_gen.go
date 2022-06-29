package commontest

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// SinkTest is data used by cmd/exercise and in testing.
type SinkTest struct {
	R              *http.Request
	ExpectedStatus int
	Name           string
}

// RouteTestRequests is data used by cmd/exercise and in testing.
type RouteTestRequests struct {
	Name  string
	Sinks []SinkTest
}

// UnsafeRouteRequests generates an unsafe request for each input
// and sink defined for this endpoint.
func UnsafeRouteRequests(r *common.Route, addr string) ([]SinkTest, error) {
	reqs := make([]SinkTest, 0, len(r.Inputs)*len(r.Sinks))
	for _, s := range r.Sinks {
		if len(s.Name) == 0 || s.Name == "_" {
			continue
		}
		for _, i := range r.Inputs {
			method := common.MethodFromInput(i)
			var u string
			if r.UsesGenericTmpl() {
				// different param order, to more easily work with gin
				u = path.Join(addr, r.Base, s.URL, i, "unsafe")
			} else {
				u = path.Join(addr, r.Base, i, s.URL, "unsafe")
			}
			// path.Join would break this by cleaning '//' to '/'
			u = "http://" + u
			req, err := http.NewRequest(method, u, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}
			s.AddPayloadToRequest(req, i, "", r.Payload)

			reqs = append(reqs, SinkTest{
				Name:           strings.Join([]string{s.Name, i}, "/"),
				R:              req,
				ExpectedStatus: s.ExpectedUnsafeStatus,
			})
		}
	}
	return reqs, nil
}

// UnsafeRequests generates a list of requests for all unsafe
// endpoints common to all app frameworks.
func UnsafeRequests(addr string, rmap common.RouteMap) ([]RouteTestRequests, error) {
	var reqs []RouteTestRequests
	if len(rmap) == 0 {
		return nil, fmt.Errorf("init error - no routes returned by GetRouteMap")
	}
	for _, route := range rmap {
		sinks, err := UnsafeRouteRequests(&route, addr)
		if err != nil {
			return nil, fmt.Errorf("route %s: %w", route.Base, err)
		}
		reqs = append(reqs, RouteTestRequests{
			Name:  route.Name,
			Sinks: sinks,
		})
	}
	return reqs, nil
}
