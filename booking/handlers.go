package main

import (
	"context"
	"log"

	bookingpb "github.com/rajsekharde/ticketing-microservices/proto/booking"
	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	paymentpb "github.com/rajsekharde/ticketing-microservices/proto/payment"

	// paymentpb "github.com/rajsekharde/ticketing-microservices/proto/payment"

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
	amount := len(eventSeatIds) * int(lockResp.TicketPrice)
	paymentReq := paymentpb.CompletePaymentRequest{
		BookingId: int64(bookingId),
		Amount: int64(amount),
	}
	paymentResp, err := cl.payment.Stub.CompletePayment(ctx, &paymentReq)
	if err != nil || paymentResp == nil {
		log.Printf("[FAILED] Book Seat, booking id: %v, Payment failed\n", bookingId)
		return nil, status.Errorf(codes.Internal, "payment failed")
	}

	// book or release seat

	log.Printf("Booking completed. Booking id: %v\n", bookingId)
	return &bookingpb.BookSeatResponse{
		BookingId: int64(bookingId),
	}, nil
}