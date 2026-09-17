package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

var user *userClient

func main() {
	port := flag.String("port", "8000", "Server port on host machine")
	flag.Parse()

	if *port == "" {
		*port = "8000"
	}


	var err error
	user, err = newUserClient("localhost:50051")
	if err != nil {
		log.Fatal(err)
	}
	defer user.close()


	router := gin.Default()

	router.GET("/health", getHealth)
	router.GET("/users/:id", getUserById)

	addr := fmt.Sprintf(":%s", *port)
	err = router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}