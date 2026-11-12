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
	
	venueId, err := db.createVenueWithSeatsQuery(v)
	if err != nil {
		log.Printf("[FAILED] Create Venue: name = %v, error: %v\n", req.Name, err)
		return nil, status.Errorf(codes.Internal, "could not create venue")
	}

	log.Printf("Create Venue: name = %v\n", req.Name)
	return &eventpb.CreateVenueResponse{VenueId: int64(venueId)}, nil
}

func (s *server) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	e := &event{
		name: req.Name,
		venueId: int(req.VenueId),
		startTime: req.StartTime.AsTime(),
		durationMinutes: int(req.DurationMinutes),
	}

	eventId, err := db.createEventQuery(e)
	if err != nil {
		log.Printf("[FAILED] Create Event: name = %v, error: %v\n", req.Name, err)
		if err.Error() == "venue not found" {
			return nil, status.Errorf(codes.NotFound, "invalid venue id")
		}
		return nil, status.Errorf(codes.Internal, "could not create event")
	}
	
	log.Printf("Create Event: name = %v\n", req.Name)
	return &eventpb.CreateEventResponse{EventId: int64(eventId)}, nil
}