package commontest

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// UnsafeRequests generates a list of requests for all unsafe
// endpoints common to all app frameworks.
func UnsafeRequests(addr string) ([]*http.Request, error) {
	var reqs []*http.Request
	type generator func(addr string) ([]*http.Request, error)
	for _, gen := range []generator{
		unsafeCmdInjections,
		unsafePathTraversals,
		unsafeSQLInjections,
		unsafeSSRF,
		unsafeXSS,
	} {
		newReqs, err := gen(addr)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, newReqs...)
	}
	return reqs, nil
}

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
			AddUserInput(t, r, "", "hello there!; echo hack hack hack")
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}

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
			file := "README.md"
			if strings.Contains(fn, "Write") {
				file = os.DevNull // We don't actually want to overwrite files
			}
			AddUserInput(t, r, "", file)
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}

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
			AddUserInput(t, r, "", "<script>alert(1);</script>")
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}

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
				AddUserInput(t, r, "", "Robert'; DROP TABLE Students;--")
			}
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}

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
			// This will make the payload match the app, but we need to fix the app.
			if t == "path" {
				t = "query"
			}
			AddUserInput(t, r, "", "www.example.com")
			reqs = append(reqs, r)
		}
	}

	return reqs, nil
}
