package servechi

import (
	"testing"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servetest"
)

func TestRouteData(t *testing.T) {
	Setup()
	t.Cleanup(func() {
		common.Reset()
	})
	servetest.TestRouteData(t, nil, nil)
}
