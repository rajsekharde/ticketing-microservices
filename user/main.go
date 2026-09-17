package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var db *database

type server struct {
	userpb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	log.Printf("Get User: id = %v\n", req.GetId())
	return &userpb.User{
		Id: req.Id,
		Email: "test@mail.com",
		Name: "RSD",
		Role: userpb.UserRole_CUSTOMER,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
    err := db.createUser(&createUserRequest{
        Email: req.Email,
        Name:  req.Name,
        Role:  req.Role.String(),
    })
    if err != nil {
        log.Printf("[FAILED] Create User: email = %v, error: %v\n", req.Email, err.Error())
        // Return a proper gRPC status error (e.g., AlreadyExists if email is taken)
        return nil, status.Errorf(codes.Internal, "failed to insert user into database")
    }

    log.Printf("Create User: email = %v\n", req.Email)
    return &userpb.CreateUserResponse{}, nil // No error field needed in response
}

func main() {
	cfg := loadEnv(".env")

	var err error
	db, err = newDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.close()
	log.Println("Connected to database")

	// create tables
	err = db.initTables()
	if err != nil {
		log.Fatalf("Failed to initialize DB tables: %v", err)
	}
	log.Println("Database tables initialized")

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
	db_host string
	db_port string
	db_user string
	db_password string
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
		db_host: os.Getenv("POSTGRES_HOST"),
		db_port: os.Getenv("POSTGRES_PORT"),
		db_user: os.Getenv("POSTGRES_USER"),
		db_password: os.Getenv("POSTGRES_PASSWORD"),
	}
}