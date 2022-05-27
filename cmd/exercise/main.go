package main

import (
	"flag"
	"log"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// defaultAddr is the default address to send requests to.
const defaultAddr = "localhost:8080"

func main() {
	log.Println("Starting...")
	addr := flag.String("addr", defaultAddr, "`host:port` to access the listening server")
	flag.Parse()
	logger := common.NewLogWrapper(log.Default())
	if err := exercise(logger, true, *addr); err != nil {
		log.Fatal(err)
	}
}
