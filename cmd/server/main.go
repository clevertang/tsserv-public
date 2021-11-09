package main

import (
	"flag"
	"log"

	"github.com/tinkermode/tsserv/pkg/tsserv"
)

func main() {
	var port int
	const defaultPort = 8080

	flag.IntVar(&port, "p", defaultPort, "port number of server")
	flag.Parse()

	log.Printf("Starting server on port %d\n", port)

	s := tsserv.New(port)

	log.Fatalf("Server exited with error %v\n", s.Run())
}
