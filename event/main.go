package main

import (
	"log"
	"net"

	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/grpc"
)

type server struct {
	eventpb.UnimplementedEventServiceServer
}

func main() {
	addr := ":50052"
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