package main

import (
	"context"
	"log"

	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) CreateVenue(ctx context.Context, req *eventpb.CreateVenueRequest) (*eventpb.CreateVenueResponse, error) {
	err := db.createVenueQuery(&venue{name: req.Name, address: req.Address})
	if err != nil {
		log.Printf("[FAILED] Create Venue: name = %v\n", req.Name)
		return nil, status.Errorf(codes.Internal, "could not create venue")
	}
	return &eventpb.CreateVenueResponse{}, nil
}