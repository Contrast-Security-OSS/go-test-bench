package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servestd"
)

// DefaultAddr is where we listen if not overridden with '-addr' flag
const DefaultAddr = "localhost:8080"

func main() {
	// set up command line flags
	flag.StringVar(&servestd.Pd.Addr, "addr", DefaultAddr, "listen on this `host:port`")
	flag.BoolVar(&common.Verbose, "v", true, "increase verbosity")
	flag.Parse()

	servestd.Setup()
	log.Fatal(http.ListenAndServe(servestd.Pd.Addr, nil))
}
