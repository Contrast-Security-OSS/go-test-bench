package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servegin"
)

// DefaultPort is the port that the API runs on if no command line argument is specified
const DefaultPort = 8080

func main() {
	// Setup command line flags
	port := flag.Int("port", DefaultPort, "listen on this `port` on localhost")
	flag.Parse()
	addr := fmt.Sprintf(":%d", *port)

	router, dbFile := servegin.Setup(addr)

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

	log.Printf("Server startup at: %s\n", addr)
	log.Fatal(router.Run(addr))
}
