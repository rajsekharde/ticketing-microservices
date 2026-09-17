package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"
	"google.golang.org/grpc"
)

type server struct {
	userpb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	log.Printf("Get User: %v\n", req.GetId())
	return &userpb.User{
		Id: req.Id,
		Email: "test@mail.com",
		Name: "RSD",
		Role: userpb.UserRole_CUSTOMER,
	}, nil
}

func main() {
	cfg := loadEnv(".env")

	addr := ":" + cfg.port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &server{})

	log.Printf("User gRPC server listening on: %v", cfg.port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

type config struct {
	port string
}

func loadEnv(path string) *config {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("Failed to load env file. Using environment variables.")
	}

	port := os.Getenv("USER_PORT")
	if port == "" {
		port = "50051"
	}

	return &config{
		port: port,
	}
}