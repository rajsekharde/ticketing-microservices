package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	notifpb "github.com/rajsekharde/ticketing-microservices/proto/notification"
	"google.golang.org/grpc"
)

type server struct {
	notifpb.UnimplementedNotificationServiceServer
}

func (s *server) GetHealth(ctx context.Context, req *notifpb.GetHealthRequest) (*notifpb.GetHealthResponse, error) {
	status := "UP"
	dbStatus := "NO_DB"
	msg := ""

	return &notifpb.GetHealthResponse{
		Status: status,
		DbStatus: dbStatus,
		Msg: msg,
	}, nil
}

var cfg config

func main() {
	cfg = loadEnv(".env")
	var err error

	addr := ":" + cfg.port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	grpcServer := grpc.NewServer()
	notifpb.RegisterNotificationServiceServer(grpcServer, &server{})
	
	log.Printf("Notification gRPC server listening on: %v\n", cfg.port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

type config struct {
	port string
}
func loadEnv(path string) config {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("Failed to load env file. Using environment variables")
	}

	port := os.Getenv("NOTIFICATION_PORT")
	if port == "" {
		port = "50055"
	}

	return config{
		port: port,
	}
}