package main

import (
	"fmt"
	"net/http"
)

func unsafeSSRF(addr string) ([]*http.Request, error) {
	types := []string{"query", "path"}
	fns := []string{"default", "http", "request"}

	reqs := make([]*http.Request, 0, len(types)*len(fns))
	for _, fn := range fns {
		for _, t := range types {
			method := http.MethodGet
			// TODO: This route is non-standard. Missing /unsafe and flipped param order for sink and source
			r, err := http.NewRequest(method, fmt.Sprintf("http://%s/ssrf/%s/%s", addr, fn, t), nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}
			// BUG: The path route is actually a query param.
			// This will make the payload match the site, but we need to fix the site.
			if t == "path" {
				t = "query"
			}
			addUserInput(t, r, "", "www.example.com")
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}
