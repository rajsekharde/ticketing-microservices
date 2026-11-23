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
	router.POST("/users/login", handlers.UserLogin)
	router.GET("/events", handlers.GetEvents)
	router.GET("/events/:id/seats", handlers.GetEventSeats)

	protected := router.Group("/")
	protected.Use(authMiddleware([]byte(cfg.jwtSecretKey), cfg.jwtIssuer))
	{
		protected.GET("/users/:id", handlers.GetUserById)
		protected.POST("/venues", handlers.CreateVenue)
		protected.POST("/events", handlers.CreateEvent)
	}

	addr := fmt.Sprintf(":%s", cfg.port)
	err := router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
