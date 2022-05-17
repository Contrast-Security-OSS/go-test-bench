package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common/commontest"
)

func exercise(log logger, addr string) error {
	client := http.DefaultClient

	// Send request to app root to determine framework
	res, err := client.Get("http://" + addr)
	if err != nil {
		return fmt.Errorf("failed to GET root: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccessful root response: %d", res.StatusCode)
	}

	// Generate requests
	var reqs []*http.Request
	f := res.Header.Get("Application-Framework")
	switch f {
	case "Stdlib", "Gin":
		reqs, err = commontest.UnsafeRequests(addr)
	case "":
		log.Fatalf("failed to determine application framework: no Application-Framework header")
	default:
		log.Fatalf("unsupported application framework: %s", f)
	}
	if err != nil {
		log.Fatalf("failed to generate requests for %s framework: %s", f, err)
	}

	// Exercise app
	for _, r := range reqs {
		if strings.Contains(r.URL.Path, "unvalidatedRedirect") {
			// temporary fixup - unvalidated redirect does not specify the input type in url
			r.URL.Path = strings.ReplaceAll(r.URL.Path, "query/", "")
		}
		log.Logf("sending request: %s", r.URL.String())
		res, err := client.Do(r)
		if err != nil {
			log.Errorf("%s: %s", r.URL, err)
		}
		expectFail := false
		if runtime.GOOS == "windows" {
			if strings.Contains(r.URL.Path, "pathTraversal/gin.File") {
				// fails on windows because /etc/passwd does not exist
				// other funcs do not report 404 for missing file
				expectFail = true
			}
		}
		reqFailed := res.StatusCode != 200
		if expectFail != reqFailed {
			log.Errorf("expected fail=%t, got fail=%t: status=%d url=%s", expectFail, reqFailed, res.StatusCode, r.URL)
		} else {
			log.Logf("route exercised: %s", r.URL.String())
		}
	}

	log.Logf("All routes exercised")
	return nil
}

// contains selected methods of testing.TB
type logger interface {
	Logf(f string, va ...interface{})
	Errorf(f string, va ...interface{})
	Fatalf(f string, va ...interface{})
}

// logWrapper wraps log.Logger to work with the above logger
// interface. Errorf() and Fatalf() are equivalent.
type logWrapper struct {
	l *log.Logger
}

func (l *logWrapper) Logf(f string, va ...interface{})   { l.l.Printf(f, va...) }
func (l *logWrapper) Errorf(f string, va ...interface{}) { l.l.Fatalf(f, va...) }
func (l *logWrapper) Fatalf(f string, va ...interface{}) { l.l.Fatalf(f, va...) }
