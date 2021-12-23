package main

import (
	"flag"
	"log"
)

const DefaultAddr = "localhost:8080"

func main() {
	log.Println("Starting...")
	addr := flag.String("addr", DefaultAddr, "`host:port` to access the listening server")
	flag.Parse()
	if err := exercise(*addr); err != nil {
		log.Fatal(err)
	}
}
