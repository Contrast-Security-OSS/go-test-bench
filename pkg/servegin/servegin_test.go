package servegin

import (
	"testing"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servetest"
)

func TestRouteData(t *testing.T) {
	Setup("don't care")
	t.Cleanup(common.Reset)
	servetest.TestRouteData(t, nil, nil)
}
