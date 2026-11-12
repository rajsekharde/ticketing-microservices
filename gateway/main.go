package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var user *userClient
var event *eventClient

func main() {
	var err error

	cfg := loadEnv(".env")

	user, err = newUserClient(cfg.userAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer user.close()

	event, err = newEventClient(cfg.eventAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer event.close()

	router := gin.Default()

	router.GET("/health", getHealth)
	router.POST("/user/register", createUser)
	router.GET("/user/:id", getUserById)
	router.POST("/user/login", userLogin)
	router.POST("/venue", createVenue)
	router.POST("/event", createEvent)

	addr := fmt.Sprintf(":%s", cfg.port)
	err = router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

type config struct {
	port string
	userAddress string
	eventAddress string
}

func loadEnv(path string) *config {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("Failed to load env file. Using environment variables.")
	}

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8000"
	}

	userHost := os.Getenv("USER_HOST")
	if userHost == "" {
		userHost = "localhost"
	}
	userPort := os.Getenv("USER_PORT")
	if userPort == "" {
		userPort = "50051"
	}
	userAddress := userHost + ":" + userPort

	eventHost := os.Getenv("EVENT_HOST")
	if eventHost == "" {
		eventHost = "localhost"
	}
	eventPort := os.Getenv("EVENT_PORT")
	if eventPort == "" {
		eventPort = "50052"
	}
	eventAddress := eventHost + ":" + eventPort

	return &config{
		port: port,
		userAddress: userAddress,
		eventAddress: eventAddress,
	}
}