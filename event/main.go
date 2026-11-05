package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/grpc"
)

type server struct {
	eventpb.UnimplementedEventServiceServer
}

func main() {
	cfg := loadEnv(".env")

	addr := ":" + cfg.port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	grpcServer := grpc.NewServer()
	eventpb.RegisterEventServiceServer(grpcServer, &server{})

	log.Printf("Event gRPC server listening on: %v\n", 50052)
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

	port := os.Getenv("EVENT_PORT")
	if port == "" {
		port = "50052"
	}

	return config{
		port: port,
	}
}