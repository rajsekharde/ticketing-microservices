package main

import (
	"context"
	"log"

	// "time"

	paymentpb "github.com/rajsekharde/ticketing-microservices/proto/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetHealth(ctx context.Context, req *paymentpb.GetHealthRequest) (*paymentpb.GetHealthResponse, error) {
	status := "UP"
	dbStatus := "CONNECTED"
	msg := ""

	err := db.db.Ping()
	if err != nil {
		dbStatus = "DISCONNECTED"
		msg = err.Error()
	}

	return &paymentpb.GetHealthResponse{
		Status: status,
		DbStatus: dbStatus,
		Msg: msg,
	}, nil
}

func (s *server) CompletePayment(ctx context.Context, req *paymentpb.CompletePaymentRequest) (*paymentpb.CompletePaymentResponse, error) {
	
	// create new transaction entry
	transId, err := db.newTransactionQuery(int(req.BookingId), int(req.Amount))
	if err != nil {
		log.Printf("[FAILED] Payment, booking id: %v, Error creating new transaction: %v\n", req.BookingId, err)
		return nil, status.Errorf(codes.Internal, "failed to create new transaction")
	}

	// payment logic
	// time.Sleep(5*time.Second)

	// update transaction status
	err = db.updateTransStatusQuery(transId, "COMPLETED")
	if err != nil {
		log.Printf("[FAILED] Payment, booking id: %v, Error updating transaction status: %v\n", req.BookingId, err)
		return nil, status.Errorf(codes.Internal, "failed to update transaction status")
	}

	return &paymentpb.CompletePaymentResponse{
		TransactionId: int64(transId),
	}, nil
}