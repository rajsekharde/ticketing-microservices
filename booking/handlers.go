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
		Status:   status,
		DbStatus: dbStatus,
		Msg:      msg,
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
		BookingId:   int64(bookingId),
	}
	lockResp, err := cl.event.LockSeat(ctx, &lockReq)
	if err != nil || lockResp == nil {
		log.Printf("[FAILED] Book Seat, booking id: %v, Failed to lock seat\n", bookingId)

		// set booking status to FAILED
		err := db.updateBookingStatusQuery(bookingId, "FAILED")
		if err != nil {
			log.Printf("Booking id: %v, error updating booking status: %v\n", bookingId, err)
		}

		return nil, status.Errorf(codes.Internal, "failed to lock event seat")
	}

	// complete payment
	amount := len(eventSeatIds) * int(lockResp.TicketPrice)
	paymentReq := paymentpb.CompletePaymentRequest{
		BookingId: int64(bookingId),
		Amount:    int64(amount),
	}
	paymentResp, err := cl.payment.Stub.CompletePayment(ctx, &paymentReq)
	if err != nil {
		log.Printf("[FAILED] Book Seat, booking id: %v, error completing payment: %v\n", bookingId, err)

		// release seat
		releaseReq := eventpb.ReleaseSeatRequest{
			EventSeatId: req.EventSeatId,
			BookingId:   int64(bookingId),
		}
		_, err := cl.event.ReleaseSeat(ctx, &releaseReq)
		if err != nil {
			log.Printf("[FAILED] Book Seat, booking id: %v, error releasing locked seat: %v\n", bookingId, err)
		} else {
			log.Printf("Book Seat, booking id: %v, Locked seat released\n", bookingId)
		}

		// set booking status to FAILED
		err = db.updateBookingStatusQuery(bookingId, "FAILED")
		if err != nil {
			log.Printf("Booking id: %v, error updating booking status: %v\n", bookingId, err)
		}

		return nil, status.Errorf(codes.Internal, "payment failed, locked seat released")
	}

	// book seat
	bookReq := eventpb.BookSeatRequest{
		EventSeatId: req.EventSeatId,
		BookingId:   int64(bookingId),
	}
	_, err = cl.event.BookSeat(ctx, &bookReq)
	if err != nil {
		log.Printf("[FAILED] Book Seat, booking id: %v, error booking locked seat: %v\n", bookingId, err)

		// set booking status to FAILED
		err = db.updateBookingStatusQuery(bookingId, "FAILED")
		if err != nil {
			log.Printf("Booking id: %v, error updating booking status: %v\n", bookingId, err)
		}

		return nil, status.Errorf(codes.Internal, "payment completed, failed to book locked seat")
	}

	// set booking status to COMPLETED
	err = db.updateBookingStatusQuery(bookingId, "COMPLETED")
	if err != nil {
		log.Printf("Booking id: %v, error updating booking status: %v\n", bookingId, err)
	}

	log.Printf("Booking completed. Booking id: %v, Transaction id: %v\n", bookingId, paymentResp.TransactionId)
	return &bookingpb.BookSeatResponse{
		BookingId: int64(bookingId),
	}, nil
}
