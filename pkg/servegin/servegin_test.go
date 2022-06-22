package servegin

import (
	"os"
	"testing"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servetest"
)

func TestRouteData(t *testing.T) {
	_,dbFile:=Setup("don't care")
	t.Cleanup(func() {
		common.Reset(		)
		_=os.RemoveAll(dbFile)
	}
	servetest.TestRouteData(t, nil, nil)
}
