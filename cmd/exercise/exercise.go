package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

type exercises struct {
	client  *http.Client
	log     common.Logger
	verbose bool
	addr    string
	reqs    []common.RouteTestRequests
}

func exercise(log common.Logger, verbose bool, addr string) error {
	e := &exercises{
		log:     log,
		verbose: verbose,
		addr:    addr,
	}
	if err := e.init(); err != nil {
		return err
	}
	// Generate requests
	for _, r := range e.reqs {
		for _, s := range r.Sinks {
			e.run(e.log, s)
		}
	}
	log.Logf("All routes exercised")
	return nil
}

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
	e.reqs, err = common.UnsafeRequests(e.addr)
	if err != nil {
		e.log.Fatalf("failed to generate requests for %s framework: %s", f, err)
	}
	return nil
}

func (e *exercises) run(log common.Logger, s common.SinkTest) {

	// Exercise app
	// for _, r := range reqs {
	if s.Status == 0 {
		s.Status = http.StatusOK
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
	if res.StatusCode != s.Status {
		log.Errorf("expected status=%d, got status=%d with url=%s", s.Status, res.StatusCode, s.R.URL)
	} else {
		if e.verbose {
			log.Logf("route exercised: %s", s.R.URL.String())
		}
	}
}
