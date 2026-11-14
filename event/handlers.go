package main

import (
	"context"
	"log"

	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) GetHealth(ctx context.Context, req *eventpb.GetHealthRequest) (*eventpb.GetHealthResponse, error) {
	status := "UP"
	dbStatus := "CONNECTED"
	msg := ""

	err := db.db.Ping()
	if err != nil {
		dbStatus = "DISCONNECTED"
		msg = err.Error()
	}

	return &eventpb.GetHealthResponse{
		Status: status,
		DbStatus: dbStatus,
		Msg: msg,
	}, nil
}

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
		ticketPrice: int(req.TicketPrice),
	}

	eventId, err := db.createEventWithSeatsQuery(e)
	if err != nil {
		log.Printf("[FAILED] Create Event: name = %v, error: %v\n", req.Name, err)
		if err.Error() == "venue not found" {
			return nil, status.Errorf(codes.NotFound, "invalid venue id")
		}
		return nil, status.Errorf(codes.Internal, "failed to create event")
	}
	
	log.Printf("Create Event: name = %v\n", req.Name)
	return &eventpb.CreateEventResponse{EventId: int64(eventId)}, nil
}

func (s *server) GetEvents(ctx context.Context, req *eventpb.GetEventsRequest) (*eventpb.GetEventsResponse, error) {
	events, err := db.getEventsQuery()
	if err != nil {
		log.Println("[FAILED] Get Events")
		return nil, status.Errorf(codes.Internal, "failed to retrieve events")
	}

	var grpcEvents []*eventpb.Event
	for _, e := range events {
		ev := &eventpb.Event{
			Id: int64(e.id),
			Name: e.name,
			VenueName: e.venueName,
			VenueAddress: e.venueAddress,
			StartTime: timestamppb.New(e.startTime),
			DurationMinutes: int64(e.durationMinutes),
			TicketPrice: int64(e.ticketPrice),
		}

		grpcEvents = append(grpcEvents, ev)
	}
	
	log.Println("Get Events")
	return &eventpb.GetEventsResponse{Events: grpcEvents}, nil
}

func (s *server) GetEventSeats(ctx context.Context, req *eventpb.GetEventSeatsRequest) (*eventpb.GetEventSeatsResponse, error) {
	eventSeats, err := db.getEventSeatsQuery(int(req.EventId))
	if err != nil {
		log.Printf("[FAILED] Get event seats: event id = %v, error: %v\n", req.EventId, err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve event seats")
	}

	var grpcEventSeats []*eventpb.Seat
	for _, e := range eventSeats {
		s := &eventpb.Seat{
			Id: int64(e.id),
			EventId: int64(e.eventId),
			SeatId: int64(e.seatId),
			Status: e.status,
			Price: int64(e.price),
		}
		grpcEventSeats = append(grpcEventSeats, s)
	}

	log.Printf("Get event seats: event id = %v\n", req.EventId)
	return &eventpb.GetEventSeatsResponse{Seats: grpcEventSeats}, nil
}