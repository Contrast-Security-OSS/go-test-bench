package main

import (
	"fmt"
	"net/http"
)

func unsafeXSS(addr string) ([]*http.Request, error) {
	types := []string{"query", "buffered-query", "params", "body", "buffered-body", "response"}
	fns := []string{"reflectedXss"}

	reqs := make([]*http.Request, 0, len(types)*len(fns))
	for _, fn := range fns {
		for _, t := range types {
			method := http.MethodGet
			if t == "headers-json" || t == "body" {
				method = http.MethodPost
			}
			r, err := http.NewRequest(method, fmt.Sprintf("http://%s/xss/%s/%s/unsafe", addr, t, fn), nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}
			addUserInput(t, r, "", "<script>alert(1);</script>")
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}
