package main

import (
	"fmt"
	"net/http"
)

func unsafeUnvalidatedRedirect(addr string) ([]*http.Request, error) {
	fns := []string{"gin.Redirect"}

	reqs := make([]*http.Request, 0, len(fns))
	for _, fn := range fns {
		method := http.MethodGet
		// TODO: This route is non-standard. Missing source param type
		r, err := http.NewRequest(method, fmt.Sprintf("http://%s/unvalidatedRedirect/%s/unsafe", addr, fn), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		addUserInput("query", r, "", "http://www.example.com")
		reqs = append(reqs, r)
	}

	return reqs, nil
}
