package commontest

import (
	"fmt"
	"net/http"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// UnsafeRequests generates a list of requests for all unsafe
// endpoints common to all app frameworks.
func UnsafeRequests(addr string) ([]*http.Request, error) {
	var reqs []*http.Request
	rmap := common.GetRouteMap()
	if len(rmap) == 0 {
		return nil, fmt.Errorf("init error - no routes returned by GetRouteMap")
	}
	for _, route := range rmap {
		r, err := route.UnsafeRequests(addr)
		if err != nil {
			return nil, fmt.Errorf("route %s: %w", route.Base, err)
		}
		reqs = append(reqs, r...)
	}
	return reqs, nil
}
