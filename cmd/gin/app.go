package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servegin"
)

// DefaultAddr is where we listen if not overridden with '-addr' flag
const DefaultAddr = "localhost:8080"

func main() {
	// set up command line flag
	addr := flag.String("addr", DefaultAddr, "listen on this `host:port`")
	flag.Parse()

	router, dbFile := servegin.Setup(*addr)

	// graceful shutdown to clean up database file
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		log.Println("Shutting down")
		err := os.Remove(dbFile)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	log.Printf("Server startup at: %s\n", *addr)
	log.Fatal(router.Run(*addr))
}
