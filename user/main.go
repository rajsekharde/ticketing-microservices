package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"
	"google.golang.org/grpc"
)

var db *database

type server struct {
	userpb.UnimplementedUserServiceServer
}

var cfg config

func main() {
	cfg = *loadEnv(".env")

	var err error
	db, err = newDatabase(&cfg)
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
	jwt_secret_key string
	access_token_expiry_mins int
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