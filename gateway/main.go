package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", "8000", "Server port on host machine")
	flag.Parse()

	if *port == "" {
		*port = "8000"
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Gateway running",
		})
	})

	addr := fmt.Sprintf(":%s", *port)
	err := router.Run(addr)
	if err != nil {
		log.Panic(err)
	}
}