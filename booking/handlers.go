package main

import (
	"context"
	"log"

	bookingpb "github.com/rajsekharde/ticketing-microservices/proto/booking"
	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetHealth(ctx context.Context, req *bookingpb.GetHealthRequest) (*bookingpb.GetHealthResponse, error) {
	status := "UP"
	dbStatus := "CONNECTED"
	msg := ""

	err := db.db.Ping()
	if err != nil {
		dbStatus = "DISCONNECTED"
		msg = err.Error()
	}

	return &bookingpb.GetHealthResponse{
		Status: status,
		DbStatus: dbStatus,
		Msg: msg,
	}, nil
}

func (s *server) BookSeat(ctx context.Context, req *bookingpb.BookSeatRequest) (*bookingpb.BookSeatResponse, error) {
	eventSeatIds := []int{int(req.EventSeatId)}
	
	// create new booking entry and booking items
	bookingId, err := db.createBookingQuery(eventSeatIds, int(req.EventId), int(req.UserId))
	if err != nil {
		log.Printf("[FAILED] Book Seat, Error creating new booking: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to create new booking")
	}
	
	// lock seat
	lockReq := eventpb.LockSeatRequest{
		EventSeatId: req.EventSeatId,
		BookingId: int64(bookingId),
	}
	lockResp, err := cl.event.LockSeat(ctx, &lockReq)
	if err != nil || lockResp == nil {
		log.Printf("[FAILED] Book Seat, booking id: %v, Failed to lock seat\n", bookingId)
		return nil, status.Errorf(codes.Internal, "failed to lock event seat")
	}

	// complete payment

	// book or release seat

	return &bookingpb.BookSeatResponse{
		BookingId: int64(bookingId),
	}, nil
}