// Command gen generates go-swagger yaml and handlers from route data.
package main

import (
	_ "embed"
	"testing"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

func Test_routePkg(t *testing.T) {

	tests := []struct {
		name, want string
		rt         common.Route
	}{
		{
			name: "basic - trims leading slash",
			rt: common.Route{
				Base: "/a",
			},
			want: "a",
		},
		{
			name: "camel case to snake case",
			rt: common.Route{
				Base: "/cmdInjection",
			},
			want: "cmd_injection",
		},
		{
			name: "camel case to snake case works for adjacent letters",
			rt: common.Route{
				Base: "/XaaaaaaaaaaajAAAAnnnnnnnnnnnnnnnnnn",
			},
			want: "xaaaaaaaaaaaj_a_a_a_annnnnnnnnnnnnnnnnn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := routePkg(&tt.rt); got != tt.want {
				t.Errorf("routePkg() = %v, want %v", got, tt.want)
			}
		})
	}
}
