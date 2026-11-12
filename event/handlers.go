package main

import (
	"context"
	"log"

	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) CreateVenue(ctx context.Context, req *eventpb.CreateVenueRequest) (*eventpb.CreateVenueResponse, error) {
	v := &venue{
		name: req.Name,
		address: req.Address,
		capacity: int(req.Capacity),
	}
	
	err := db.createVenueQuery(v)
	if err != nil {
		log.Printf("[FAILED] Create Venue: name = %v, error: %v\n", req.Name, err)
		return nil, status.Errorf(codes.Internal, "could not create venue")
	}

	err = db.createSeatsQuery(v)
	if err != nil {
		log.Printf("[FAILED] Create Seats: venue name = %v, error: %v\n", req.Name, err)
		return nil, status.Errorf(codes.Internal, "could not create seats")
	}

	log.Printf("Create Venue: name = %v\n", req.Name)
	return &eventpb.CreateVenueResponse{}, nil
}