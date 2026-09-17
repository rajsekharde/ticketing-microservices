package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", "8000", "Server port on host machine")
	flag.Parse()

	if *port == "" {
		*port = "8000"
	}

	router := gin.Default()

	router.GET("/health", getHealth)

	addr := fmt.Sprintf(":%s", *port)
	err := router.Run(addr)
	if err != nil {
		log.Panic(err)
	}
}