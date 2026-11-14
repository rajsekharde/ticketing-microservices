package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"
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

	router.GET("/health", getHealthHandler)
	router.POST("/user/register", createUser)
	router.GET("/user/:id", getUserById)
	router.POST("/user/login", userLogin)
	router.POST("/venue", createVenue)
	router.POST("/event", createEvent)
	router.GET("/events", getEvents)
	router.GET("/event/:id/seats", getEventSeats)

	addr := fmt.Sprintf(":%s", cfg.port)
	err = router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}

type config struct {
	port         string
	userAddress  string
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
		port:         port,
		userAddress:  userAddress,
		eventAddress: eventAddress,
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
	var userHealth, eventHealth serviceHealth

	start := time.Now()
	r1, err := user.stub.GetHealth(c.Request.Context(), &userpb.GetHealthRequest{})
	if err != nil || r1 == nil {
		userHealth = serviceHealth{
			Name:      "user",
			Status:    "DOWN",
			DbStatus:  "DISCONNECTED",
			Msg:       err.Error(),
			LatencyMs: float32(time.Since(start).Milliseconds()),
		}
	} else {
		userHealth = serviceHealth{
			Name:      "user",
			Status:    r1.Status,
			DbStatus:  r1.DbStatus,
			Msg:       r1.Msg,
			LatencyMs: float32(time.Since(start).Milliseconds()),
		}
	}
	health = append(health, userHealth)

	start = time.Now()
	r2, err := event.stub.GetHealth(c.Request.Context(), &eventpb.GetHealthRequest{})
	if err != nil || r2 == nil {
		eventHealth = serviceHealth{
			Name:      "event",
			Status:    "DOWN",
			DbStatus:  "DISCONNECTED",
			Msg:       err.Error(),
			LatencyMs: float32(time.Since(start).Milliseconds()),
		}
	} else {
		eventHealth = serviceHealth{
			Name:      "event",
			Status:    r2.Status,
			DbStatus:  r2.DbStatus,
			Msg:       r2.Msg,
			LatencyMs: float32(time.Since(start).Milliseconds()),
		}
	}
	health = append(health, eventHealth)

	resp := healthStruct{
		RecordedAt:    time.Now(),
		GatewayStatus: "RUNNING",
		GrpcServices:  health,
	}
	c.IndentedJSON(http.StatusOK, resp)
}
