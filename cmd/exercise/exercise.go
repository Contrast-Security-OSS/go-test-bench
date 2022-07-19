package main

import (
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/serveswagger"
	"io"
	"net/http"
	"strings"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common/commontest"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servegin"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servestd"
)

type exercises struct {
	client     *http.Client
	verbose    bool
	standalone bool
	addr       string
	reqs       []commontest.RouteTestRequests
	framework  string
	rmap       common.RouteMap
}

// exercise was split up so that tests can report sub-test names
func exercise(log common.Logger, verbose bool, addr string) error {
	e := &exercises{
		verbose:    verbose,
		addr:       addr,
		standalone: true,
	}
	// create requests
	e.init(log)
	e.checkAssets(log)

	// send requests
	for _, r := range e.reqs {
		for _, s := range r.Sinks {
			e.run(log, s)
		}
	}
	log.Logf("All routes exercised")
	return nil
}

// determine framework, then create (but do not send) requests
func (e *exercises) init(log common.Logger) {
	e.client = http.DefaultClient

	e.checkFramework(log)

	var err error
	e.reqs, err = commontest.UnsafeRequests(e.addr, e.rmap)
	if err != nil {
		log.Fatalf("failed to generate requests for %s framework: %s", e.framework, err)
	}
}

const hdrname = "Application-Framework"

// Send request to app root to identify the framework in use
func (e *exercises) checkFramework(log common.Logger) {
	res, err := e.client.Get("http://" + e.addr)
	if err != nil {
		log.Fatalf("failed to GET root: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("unsuccessful root response: %d", res.StatusCode)
	}
	e.framework = res.Header.Get(hdrname)
	switch e.framework {
	case "Stdlib":
		if e.standalone {
			servestd.RegisterRoutes()
			e.rmap = common.PopulateRouteMap(common.AllRoutes)
		}
	case "Gin":
		if e.standalone {
			servegin.RegisterRoutes()
			e.rmap = common.PopulateRouteMap(common.AllRoutes)
		}
	case "Go-Swagger":
		if e.standalone {
			serveswagger.Setup()
			routes := common.Routes{}
			swaggerRoutes := serveswagger.SwaggerParams.Rulebar

			for _, route := range swaggerRoutes {
				routes = append(routes, route)
			}

			e.rmap = common.PopulateRouteMap(routes)
		}

	case "":
		log.Fatalf("failed to determine application framework: no %q header", hdrname)
	default:
		log.Fatalf("unsupported application framework: %s", e.framework)
	}
	if e.rmap == nil {
		e.rmap = common.GetRouteMap()
	}
}

// ensure assets (currently only app.css) are loadable
func (e *exercises) checkAssets(log common.Logger) {
	res, err := e.client.Get("http://" + e.addr + "/assets/app.css")
	if err != nil {
		log.Errorf("failed to GET %s: %s", res.Request.URL, err)
	}
	wantCT := "text/css"
	if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, wantCT) {
		log.Errorf("expected content type %q, got %q", wantCT, ct)
	}
	if b, err := io.ReadAll(res.Body); err != nil || len(b) < 1024 {
		estr := "(nil)"
		if err != nil {
			estr = err.Error()
		}
		log.Errorf("undersize or unreadable css: len=%d err=%s", len(b), estr)
	}
}

// send requests
func (e *exercises) run(log common.Logger, s commontest.SinkTest) {
	if s.ExpectedStatus == 0 {
		s.ExpectedStatus = http.StatusOK
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
