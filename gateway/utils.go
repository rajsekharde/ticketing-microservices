package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	bookingpb "github.com/rajsekharde/ticketing-microservices/proto/booking"
	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	notifpb "github.com/rajsekharde/ticketing-microservices/proto/notification"
	paymentpb "github.com/rajsekharde/ticketing-microservices/proto/payment"
	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"
)

// Creates grpc clients for services
func createClients(cfg *config) {
	var err error

	user, err = newUserClient(cfg.userAddress)
	if err != nil {
		log.Println(err.Error())
	}

	event, err = newEventClient(cfg.eventAddress)
	if err != nil {
		log.Println(err.Error())
	}

	booking, err = newBookingClient(cfg.bookingAddress)
	if err != nil {
		log.Println(err.Error())
	}

	payment, err = newPaymentClient(cfg.paymentAddress)
	if err != nil {
		log.Println(err.Error())
	}

	notification, err = newNotificationClient(cfg.notificationAddress)
	if err != nil {
		log.Println(err.Error())
	}
}


type config struct {
	port                string
	userAddress         string
	eventAddress        string
	bookingAddress      string
	paymentAddress      string
	notificationAddress string
}
// Loads the contenst of a .env file as environment variables
func loadEnv(path string) *config {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("Failed to load env file. Using environment variables.")
	}

	// Gateway
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8000"
	}

	// User
	userHost := os.Getenv("USER_HOST")
	if userHost == "" {
		userHost = "localhost"
	}
	userPort := os.Getenv("USER_PORT")
	if userPort == "" {
		userPort = "50051"
	}
	userAddress := userHost + ":" + userPort

	// Event
	eventHost := os.Getenv("EVENT_HOST")
	if eventHost == "" {
		eventHost = "localhost"
	}
	eventPort := os.Getenv("EVENT_PORT")
	if eventPort == "" {
		eventPort = "50052"
	}
	eventAddress := eventHost + ":" + eventPort

	// Booking
	bookingHost := os.Getenv("BOOKING_HOST")
	if bookingHost == "" {
		bookingHost = "localhost"
	}
	bookingPort := os.Getenv("BOOKING_PORT")
	if bookingPort == "" {
		bookingPort = "50053"
	}
	bookingAddress := bookingHost + ":" + bookingPort

	// Payment
	paymentHost := os.Getenv("PAYMENT_HOST")
	if paymentHost == "" {
		paymentHost = "localhost"
	}
	paymentPort := os.Getenv("PAYMENT_PORT")
	if paymentPort == "" {
		paymentPort = "50054"
	}
	paymentAddress := paymentHost + ":" + paymentPort

	// Notification
	notificationtHost := os.Getenv("NOTIFICATION_HOST")
	if notificationtHost == "" {
		notificationtHost = "localhost"
	}
	notificationPort := os.Getenv("NOTIFICATION_PORT")
	if notificationPort == "" {
		notificationPort = "50055"
	}
	notificationAddress := notificationtHost + ":" + notificationPort

	return &config{
		port:                port,
		userAddress:         userAddress,
		eventAddress:        eventAddress,
		bookingAddress:      bookingAddress,
		paymentAddress:      paymentAddress,
		notificationAddress: notificationAddress,
	}
}

type serviceHealth struct {
	Name      string  `json:"service"`
	Status    string  `json:"status"`
	DbStatus  string  `json:"db_status"`
	Msg       string  `json:"msg"`
	LatencyMs float32 `json:"latency_ms"`
}
type healthStruct struct {
	RecordedAt    time.Time       `json:"recorded_at"`
	GatewayStatus string          `json:"gateway_status"`
	GrpcServices  []serviceHealth `json:"grpc_services"`
}
// Returns running status of all services
func getHealthHandler(c *gin.Context) {
	var health []serviceHealth

	healthChecks := []struct {
		name    string
		checkFn func(ctx context.Context) (string, string, string, error)
	}{
		{"user", func(ctx context.Context) (string, string, string, error) {
			if user == nil || user.stub == nil {
				return "", "", "", fmt.Errorf("user client not initialized")
			}
			resp, err := user.stub.GetHealth(c.Request.Context(), &userpb.GetHealthRequest{})
			if err != nil || resp == nil {
				return "", "", "", err
			}
			return resp.Status, resp.DbStatus, resp.Msg, nil
		}},
		{"event", func(ctx context.Context) (string, string, string, error) {
			if event == nil || event.stub == nil {
				return "", "", "", fmt.Errorf("client not initialized")
			}
			resp, err := event.stub.GetHealth(c.Request.Context(), &eventpb.GetHealthRequest{})
			if err != nil || resp == nil {
				return "", "", "", err
			}
			return resp.Status, resp.DbStatus, resp.Msg, nil
		}},
		{"booking", func(ctx context.Context) (string, string, string, error) {
			if booking == nil || booking.stub == nil {
				return "", "", "", fmt.Errorf("booking client not initialized")
			}
			resp, err := booking.stub.GetHealth(c.Request.Context(), &bookingpb.GetHealthRequest{})
			if err != nil || resp == nil {
				return "", "", "", err
			}
			return resp.Status, resp.DbStatus, resp.Msg, nil
		}},
		{"payment", func(ctx context.Context) (string, string, string, error) {
			if payment == nil || payment.stub == nil {
				return "", "", "", fmt.Errorf("payment client not initialized")
			}
			resp, err := payment.stub.GetHealth(c.Request.Context(), &paymentpb.GetHealthRequest{})
			if err != nil || resp == nil {
				return "", "", "", err
			}
			return resp.Status, resp.DbStatus, resp.Msg, nil
		}},
		{"notification", func(ctx context.Context) (string, string, string, error) {
			if notification == nil || notification.stub == nil {
				return "", "", "", fmt.Errorf("notification client not initialized")
			}
			resp, err := notification.stub.GetHealth(c.Request.Context(), &notifpb.GetHealthRequest{})
			if err != nil || resp == nil {
				return "", "", "", err
			}
			return resp.Status, resp.DbStatus, resp.Msg, nil
		}},
	}

	for _, service := range healthChecks {
		currHealth := serviceHealth{}
		start := time.Now()
		status, dbStatus, msg, err := service.checkFn(c.Request.Context())
		if err != nil {
			currHealth = serviceHealth{
				Name:      service.name,
				Status:    "DOWN",
				DbStatus:  "DISCONNECTED",
				Msg:       err.Error(),
				LatencyMs: float32(time.Since(start).Milliseconds()),
			}
		} else {
			currHealth = serviceHealth{
				Name:      service.name,
				Status:    status,
				DbStatus:  dbStatus,
				Msg:       msg,
				LatencyMs: float32(time.Since(start).Milliseconds()),
			}
		}
		health = append(health, currHealth)
	}

	resp := healthStruct{
		RecordedAt:    time.Now(),
		GatewayStatus: "RUNNING",
		GrpcServices:  health,
	}
	c.IndentedJSON(http.StatusOK, resp)
}
