package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

var user *userClient
var event *eventClient
var booking *bookingClient
var payment *paymentClient
var notification *notificationClient

func main() {
	cfg := loadEnv(".env")

	createClients(cfg)
	defer func() {
		user.close()
		event.close()
		booking.close()
		payment.close()
		notification.close()
	}()

	router := gin.Default()

	router.GET("/health", getHealthHandler)
	router.POST("/users/register", createUser)
	router.GET("/users/:id", getUserById)
	router.POST("/users/login", userLogin)
	router.POST("/venues", createVenue)
	router.POST("/events", createEvent)
	router.GET("/events", getEvents)
	router.GET("/events/:id/seats", getEventSeats)

	addr := fmt.Sprintf(":%s", cfg.port)
	err := router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

