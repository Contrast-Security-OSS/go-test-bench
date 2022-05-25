package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common/commontest"
)

type exercises struct {
	client  *http.Client
	log     common.Logger
	verbose bool
	addr    string
	reqs    []commontest.RouteTestRequests
}

// exercise was split up so that tests can report sub-test names
func exercise(log common.Logger, verbose bool, addr string) error {
	e := &exercises{
		log:     log,
		verbose: verbose,
		addr:    addr,
	}
	// create requests
	if err := e.init(); err != nil {
		return err
	}
	// send requests
	for _, r := range e.reqs {
		for _, s := range r.Sinks {
			e.run(e.log, s)
		}
	}
	log.Logf("All routes exercised")
	return nil
}

// determine framework, then create (but do not send) requests
func (e *exercises) init() error {
	e.client = http.DefaultClient

	// Send request to app root to determine framework
	res, err := e.client.Get("http://" + e.addr)
	if err != nil {
		return fmt.Errorf("failed to GET root: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccessful root response: %d", res.StatusCode)
	}

	f := res.Header.Get("Application-Framework")
	switch f {
	case "Stdlib", "Gin":
		//below
	case "":
		e.log.Fatalf("failed to determine application framework: no Application-Framework header")
	default:
		e.log.Fatalf("unsupported application framework: %s", f)
	}
	e.reqs, err = commontest.UnsafeRequests(e.addr)
	if err != nil {
		e.log.Fatalf("failed to generate requests for %s framework: %s", f, err)
	}
	return nil
}

// send requests
func (e *exercises) run(log common.Logger, s commontest.SinkTest) {
	if s.ExpectedStatus == 0 {
		s.ExpectedStatus = http.StatusOK
	}
	if strings.Contains(s.R.URL.Path, "unvalidatedRedirect") {
		// temporary fixup - unvalidated redirect does not specify the input type in url
		s.R.URL.Path = strings.ReplaceAll(s.R.URL.Path, "query/", "")
	}
	if e.verbose {
		log.Logf("sending request: %s", s.R.URL.String())
	}
	res, err := e.client.Do(s.R)
	if err != nil {
		log.Errorf("%s: %s", s.R.URL, err)
	}
	if res.StatusCode != s.ExpectedStatus {
		log.Errorf("expected status=%d, got status=%d with url=%s", s.ExpectedStatus, res.StatusCode, s.R.URL)
	} else {
		if e.verbose {
			log.Logf("route exercised: %s", s.R.URL.String())
		}
	}
}
