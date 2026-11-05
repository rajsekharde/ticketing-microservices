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

var cfg config

func main() {
	cfg = loadEnv(".env")

	// connect to database
	db, err := newDatabase(&cfg)
	if err != nil {
		log.Printf("Failed to connect to database: %v\n", err)
	} else {
		log.Printf("Connected to database: %v\n", cfg.db_name)
	}
	defer db.close()

	// create tables
	err = db.initTables()
	if err != nil {
		log.Printf("Failed to initialize DB tables: %v\n", err)
	} else {
		log.Println("Initialized DB tables")
	}

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
	db_host string
	db_port string
	db_user string
	db_password string
	db_name string
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
		db_host: os.Getenv("POSTGRES_HOST"),
		db_port: os.Getenv("POSTGRES_PORT"),
		db_user: os.Getenv("POSTGRES_USER"),
		db_password: os.Getenv("POSTGRES_PASSWORD"),
		db_name: os.Getenv("EVENT_DB"),
	}
}