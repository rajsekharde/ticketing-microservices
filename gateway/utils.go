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

	"github.com/rajsekharde/ticketing-microservices/gateway/clients"
	"github.com/rajsekharde/ticketing-microservices/gateway/handlers"
)

// Creates grpc clients for services and sends clients to handlers
func createClients(cfg *config) {
	var err error

	user, err = clients.NewUserClient(cfg.userAddress)
	if err != nil {
		log.Println(err.Error())
	}

	event, err = clients.NewEventClient(cfg.eventAddress)
	if err != nil {
		log.Println(err.Error())
	}

	booking, err = clients.NewBookingClient(cfg.bookingAddress)
	if err != nil {
		log.Println(err.Error())
	}

	payment, err = clients.NewPaymentClient(cfg.paymentAddress)
	if err != nil {
		log.Println(err.Error())
	}

	notification, err = clients.NewNotificationClient(cfg.notificationAddress)
	if err != nil {
		log.Println(err.Error())
	}

	// send clients to handlers package
	handlers.LoadClients(&handlers.Clients{
		User: user,
		Event: event,
		Booking: booking,
		Payment: payment,
		Notification: notification,
	})
}


type config struct {
	port                string
	userAddress         string
	eventAddress        string
	bookingAddress      string
	paymentAddress      string
	notificationAddress string
	jwtSecretKey string
	jwtIssuer string
}
// Loads the contents of a .env file as environment variables
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

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	jwtIssuer := os.Getenv("JWT_ISSUER")

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
		jwtSecretKey: jwtSecret,
		jwtIssuer: jwtIssuer,
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
			if user == nil || user.Stub == nil {
				return "", "", "", fmt.Errorf("user client not initialized")
			}
			resp, err := user.Stub.GetHealth(c.Request.Context(), &userpb.GetHealthRequest{})
			if err != nil || resp == nil {
				return "", "", "", err
			}
			return resp.Status, resp.DbStatus, resp.Msg, nil
		}},
		{"event", func(ctx context.Context) (string, string, string, error) {
			if event == nil || event.Stub == nil {
				return "", "", "", fmt.Errorf("client not initialized")
			}
			resp, err := event.Stub.GetHealth(c.Request.Context(), &eventpb.GetHealthRequest{})
			if err != nil || resp == nil {
				return "", "", "", err
			}
			return resp.Status, resp.DbStatus, resp.Msg, nil
		}},
		{"booking", func(ctx context.Context) (string, string, string, error) {
			if booking == nil || booking.Stub == nil {
				return "", "", "", fmt.Errorf("booking client not initialized")
			}
			resp, err := booking.Stub.GetHealth(c.Request.Context(), &bookingpb.GetHealthRequest{})
			if err != nil || resp == nil {
				return "", "", "", err
			}
			return resp.Status, resp.DbStatus, resp.Msg, nil
		}},
		{"payment", func(ctx context.Context) (string, string, string, error) {
			if payment == nil || payment.Stub == nil {
				return "", "", "", fmt.Errorf("payment client not initialized")
			}
			resp, err := payment.Stub.GetHealth(c.Request.Context(), &paymentpb.GetHealthRequest{})
			if err != nil || resp == nil {
				return "", "", "", err
			}
			return resp.Status, resp.DbStatus, resp.Msg, nil
		}},
		{"notification", func(ctx context.Context) (string, string, string, error) {
			if notification == nil || notification.Stub == nil {
				return "", "", "", fmt.Errorf("notification client not initialized")
			}
			resp, err := notification.Stub.GetHealth(c.Request.Context(), &notifpb.GetHealthRequest{})
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
