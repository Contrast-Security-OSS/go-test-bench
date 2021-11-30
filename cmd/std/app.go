package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servestd"
)

// DefaultPort is the port that the API runs on if no command line argument is specified
const DefaultPort = 8080

func main() {
	// set up command line flag
	port := flag.Int("port", DefaultPort, "listen on this `port` on localhost")
	flag.BoolVar(&common.Verbose, "v", true, "increase verbosity")
	flag.Parse()
	servestd.Pd.Addr = fmt.Sprintf(":%d", *port)

	servestd.Setup()
	log.Fatal(http.ListenAndServe(servestd.Pd.Addr, nil))
}
