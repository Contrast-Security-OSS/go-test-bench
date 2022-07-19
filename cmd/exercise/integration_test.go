package main

import (
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servegin"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servestd"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/serveswagger"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
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
		"Go-Swagger": {
			setup: func(t *testing.T) http.Handler {
				server, _ := serveswagger.Setup(true)

				return server.GetHandler()
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(func() {
				// TODO: Refactor to remove this global state
				common.Reset()
				http.DefaultServeMux = http.NewServeMux()
			})
			handler := test.setup(t)
			srv := httptest.NewServer(handler)
			t.Cleanup(srv.Close)
			addr := strings.TrimPrefix(srv.URL, "http://")

			e := &exercises{
				addr: addr,
			}
			e.init(t)

			t.Run("check css asset content type", func(t *testing.T) {
				e.checkAssets(t)
			})
			for _, r := range e.reqs {
				if len(r.Sinks) == 0 {
					continue
				}
				t.Run(r.Name, func(t *testing.T) {
					for _, s := range r.Sinks {
						t.Run(s.Name, func(t *testing.T) {
							e.run(t, s)
						})
					}
				})
			}
		})
	}
}
