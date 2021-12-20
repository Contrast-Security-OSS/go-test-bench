package main

import (
	"fmt"
	"net/http"
)

func unsafePathTraversals(addr string) ([]*http.Request, error) {
	types := []string{"query", "buffered-query", "headers", "body"}
	fns := []string{"ioutil.ReadFile", "ioutil.WriteFile"}

	reqs := make([]*http.Request, 0, len(types)*len(fns))
	for _, fn := range fns {
		for _, t := range types {
			method := http.MethodGet
			if t == "body" {
				method = http.MethodPost
			}
			r, err := http.NewRequest(method, fmt.Sprintf("http://%s/pathTraversal/%s/%s/unsafe", addr, t, fn), nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}
			addUserInput(t, r, "", "README.md")
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}
