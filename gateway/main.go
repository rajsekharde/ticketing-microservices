package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var user *userClient

func main() {
	var err error

	cfg := loadEnv(".env")

	user, err = newUserClient(cfg.userAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer user.close()


	router := gin.Default()

	router.GET("/health", getHealth)
	router.POST("/users", createUser)
	router.GET("/users/:id", getUserById)

	addr := fmt.Sprintf(":%s", cfg.port)
	err = router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

type config struct {
	port string
	userAddress string
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

	return &config{
		port: port,
		userAddress: userAddress,
	}
}