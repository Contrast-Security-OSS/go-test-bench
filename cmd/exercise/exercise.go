package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servegin/gintest"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servestd/stdtest"
)

func exercise(addr string) error {
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
	case "Stdlib":
		reqs, err = stdtest.UnsafeStdlibRequests(addr)
	case "Gin":
		reqs, err = gintest.UnsafeGinRequests(addr)
	case "":
		return fmt.Errorf("failed to determine application framework: no Application-Framework header")
	default:
		return fmt.Errorf("unsupported application framework: %s", f)
	}
	if err != nil {
		return fmt.Errorf("failed to generate requests for %s framework: %s", f, err)
	}

	// Exercise app
	for _, r := range reqs {
		log.Printf("sending request: %s", r.URL.String())
		res, err := client.Do(r)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return fmt.Errorf("unsuccessful response: %d", res.StatusCode)
		}
		log.Printf("route exercised: %s", r.URL.String())
	}

	log.Println("All routes successfully exercised")
	return nil
}
