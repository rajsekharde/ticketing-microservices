package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/rajsekharde/ticketing-microservices/gateway/clients"
	"github.com/rajsekharde/ticketing-microservices/gateway/handlers"
)

var user *clients.UserClient
var event *clients.EventClient
var booking *clients.BookingClient
var payment *clients.PaymentClient
var notification *clients.NotificationClient

func main() {
	cfg := loadEnv(".env")

	createClients(cfg)
	defer func() {
		user.Close()
		event.Close()
		booking.Close()
		payment.Close()
		notification.Close()
	}()

	router := gin.Default()

	router.GET("/health", getHealthHandler)
	router.POST("/users/register", handlers.CreateUser)
	router.GET("/users/:id", handlers.GetUserById)
	router.POST("/users/login", handlers.UserLogin)
	router.POST("/venues", handlers.CreateVenue)
	router.POST("/events", handlers.CreateEvent)
	router.GET("/events", handlers.GetEvents)
	router.GET("/events/:id/seats", handlers.GetEventSeats)

	addr := fmt.Sprintf(":%s", cfg.port)
	err := router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

