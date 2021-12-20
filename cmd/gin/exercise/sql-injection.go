package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func unsafeSQLInjections(addr string) ([]*http.Request, error) {
	types := []string{"headers-json", "query", "body"}
	fns := []string{"sqlite3Exec"}

	reqs := make([]*http.Request, 0, len(types)*len(fns))
	for _, fn := range fns {
		for _, t := range types {
			method := http.MethodGet
			if t == "headers-json" || t == "body" {
				method = http.MethodPost
			}
			r, err := http.NewRequest(method, fmt.Sprintf("http://%s/sqlInjection/%s/%s/unsafe", addr, t, fn), nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}
			if t == "headers-json" {
				addHeadersJSONCreds(r)
			} else {
				addUserInput(t, r, "", "Robert'; DROP TABLE Students;--")
			}
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}

func addHeadersJSONCreds(r *http.Request) {
	// TODO(james): add note about this being non-standard and tightly coupled
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	creds.Username = "Robert'; DROP TABLE Students;--"
	data, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("failed to marshal JSON object: %s", err)
	}
	r.Header.Set("credentials", string(data))
}
