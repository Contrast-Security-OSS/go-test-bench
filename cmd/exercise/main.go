package main

import (
	"flag"
	"log"
)

// defaultAddr is the default address to send requests to.
const defaultAddr = "localhost:8080"

func main() {
	log.Println("Starting...")
	addr := flag.String("addr", defaultAddr, "`host:port` to access the listening server")
	flag.Parse()
	if err := exercise(*addr); err != nil {
		log.Fatal(err)
	}
}
