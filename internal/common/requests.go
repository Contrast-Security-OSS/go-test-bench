package common

import (
	"fmt"
	"net/http"
	"strings"
)

type SinkTest struct {
	R      *http.Request
	Status int
	Name   string
}
type RouteTestRequests struct {
	Name  string
	Sinks []SinkTest
}

// UnsafeRequests generates an unsafe request for each input and sink defined
// for this endpoint.
func (r *Route) UnsafeRequests(addr string) ([]SinkTest, error) {
	reqs := make([]SinkTest, 0, len(r.Inputs)*len(r.Sinks))
	for _, s := range r.Sinks {
		if len(s.Name) == 0 || s.Name == "_" {
			continue
		}
		for _, i := range r.Inputs {
			method := methodFromInput(i)
			var u string
			if r.genericTmpl {
				// different param order, to more easily work with gin
				u = fmt.Sprintf("http://%s%s/%s/%s/unsafe", addr, r.Base, s.Name, i)
			} else {
				u = fmt.Sprintf("http://%s%s/%s/%s/unsafe", addr, r.Base, i, s.Name)
			}
			req, err := http.NewRequest(method, u, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}
			s.AddPayloadToRequest(req, i, "", r.Payload)

			reqs = append(reqs, SinkTest{
				Name:   strings.Join([]string{s.Name, i}, "/"),
				R:      req,
				Status: s.ExpectedUnsafeStatus,
			})
		}
	}
	return reqs, nil
}

// UnsafeRequests generates a list of requests for all unsafe
// endpoints common to all app frameworks.
func UnsafeRequests(addr string) ([]RouteTestRequests, error) {
	var reqs []RouteTestRequests
	rmap := GetRouteMap()
	if len(rmap) == 0 {
		return nil, fmt.Errorf("init error - no routes returned by GetRouteMap")
	}
	for _, route := range rmap {
		sinks, err := route.UnsafeRequests(addr)
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
