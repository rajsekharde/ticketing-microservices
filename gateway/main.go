package main

import (
	"flag"
	"log"
)

func main() {
	port := flag.String("port", "8000", "Server port on host machine")
	flag.Parse()

	if *port == "" {
		*port = "8000"
	}

	log.Printf("API Gateway running on port %v\n", *port)
}