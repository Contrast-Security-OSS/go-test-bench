package gintest

import (
	"fmt"
	"net/http"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common/commontest"
)

// UnsafeGinRequests generates a list of requests for all unsafe
// endpoints for the Gin version of go-test-bench.
func UnsafeGinRequests(addr string) ([]*http.Request, error) {
	common, err := commontest.UnsafeRequests(addr)
	if err != nil {
		return nil, err
	}
	ginReqs, err := unsafeGinUnvalidatedRedirect(addr)
	if err != nil {
		return nil, err
	}
	return append(common, ginReqs...), nil
}

func unsafeGinUnvalidatedRedirect(addr string) ([]*http.Request, error) {
	fns := []string{"gin.Redirect"}

	reqs := make([]*http.Request, 0, len(fns))
	for _, fn := range fns {
		method := http.MethodGet
		// TODO: This route is non-standard. Missing source param type
		r, err := http.NewRequest(method, fmt.Sprintf("http://%s/unvalidatedRedirect/%s/unsafe", addr, fn), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		commontest.AddUserInput("query", r, "", "http://www.example.com")
		reqs = append(reqs, r)
	}

	return reqs, nil
}
