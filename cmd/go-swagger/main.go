package main

import (
	"log"

	"github.com/Contrast-Security-OSS/go-test-bench/pkg/serveswagger"
)

func main() {
	server, err := serveswagger.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	defer server.Shutdown()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
