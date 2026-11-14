package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	bookingpb "github.com/rajsekharde/ticketing-microservices/proto/booking"
	"google.golang.org/grpc"
)

type server struct {
	bookingpb.UnimplementedBookingServiceServer
}

var cfg config
var db *database

func main() {
	cfg = loadEnv(".env")
	var err error

	// connect to database
	db, err = newDatabase(&cfg)
	if err != nil {
		log.Printf("Failed to connect to database: %v\n", err)
	} else {
		log.Printf("Connected to database: %v\n", cfg.db_name)
	}
	defer db.close()

	addr := ":" + cfg.port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	grpcServer := grpc.NewServer()
	bookingpb.RegisterBookingServiceServer(grpcServer, &server{})
	
	log.Printf("Booking gRPC server listening on: %v\n", cfg.port)
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

	port := os.Getenv("BOOKING_PORT")
	if port == "" {
		port = "50053"
	}

	return config{
		port: port,
		db_host: os.Getenv("POSTGRES_HOST"),
		db_port: os.Getenv("POSTGRES_PORT"),
		db_user: os.Getenv("POSTGRES_USER"),
		db_password: os.Getenv("POSTGRES_PASSWORD"),
		db_name: os.Getenv("BOOKING_DB"),
	}
}