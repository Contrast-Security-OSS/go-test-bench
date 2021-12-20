package main

import (
	"fmt"
	"net/http"
)

func unsafeCmdInjections(addr string) ([]*http.Request, error) {
	types := []string{"query", "cookies"}
	fns := []string{"exec.Command", "exec.CommandContext"}

	reqs := make([]*http.Request, 0, len(types)*len(fns))
	for _, fn := range fns {
		for _, t := range types {
			method := http.MethodGet
			if t == "cookies" {
				method = http.MethodPost
			}
			// TODO: This route is non-standard. Flipped param order for source and sink.
			r, err := http.NewRequest(method, fmt.Sprintf("http://%s/cmdInjection/%s/%s/unsafe", addr, fn, t), nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}
			addUserInput(t, r, "", "hello there!; echo hack hack hack")
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}
