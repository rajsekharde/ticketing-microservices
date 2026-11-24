package main

import (
	"log"
	"net"

	bookingpb "github.com/rajsekharde/ticketing-microservices/proto/booking"
	"google.golang.org/grpc"
)

type server struct {
	bookingpb.UnimplementedBookingServiceServer
}

var cfg config
var db *database
var cl *grpcClients

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

	// create tables
	err = db.initTables()
	if err != nil {
		log.Printf("Failed to initialize DB tables: %v\n", err)
	} else {
		log.Println("Initialized DB tables")
	}

	// create grpc clients
	cl, err = createClients(&cfg)
	if err != nil {
		log.Printf("Failed to create clients, error: %v\n", err)
	}

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