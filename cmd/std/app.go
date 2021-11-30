package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servestd"
)

// DefaultPort is the port that the API runs on if not overridden with '-addr' flag
const DefaultPort = "localhost:8080"

func main() {
	// set up command line flags
	flag.StringVar(&servestd.Pd.Addr, "addr", DefaultPort, "listen on this `host:port`")
	flag.BoolVar(&common.Verbose, "v", true, "increase verbosity")
	flag.Parse()

	servestd.Setup()
	log.Fatal(http.ListenAndServe(servestd.Pd.Addr, nil))
}
