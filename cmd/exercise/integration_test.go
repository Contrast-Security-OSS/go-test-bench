package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servegin"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servestd"
)

func TestExerciseIntegration(t *testing.T) {
	tests := map[string]struct {
		setup func(t *testing.T) http.Handler
	}{
		"Stdlib": {
			setup: func(_ *testing.T) http.Handler {
				servestd.Setup()
				return http.DefaultServeMux
			},
		},
		"Gin": {
			setup: func(t *testing.T) http.Handler {
				router, dbFile := servegin.Setup("don't care")
				t.Cleanup(func() {
					os.Remove(dbFile)
				})
				return router
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(func() {
				// TODO: Refactor to remove this global state
				common.AllRoutes = nil
				http.DefaultServeMux = http.NewServeMux()
			})
			handler := test.setup(t)
			srv := httptest.NewServer(handler)
			t.Cleanup(srv.Close)
			addr := strings.TrimPrefix(srv.URL, "http://")

			err := exercise(addr)
			if err != nil {
				t.Fatal(err)
			}
		})

	}
}
