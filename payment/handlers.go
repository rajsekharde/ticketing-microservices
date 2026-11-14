package main

import (
	"context"

	paymentpb "github.com/rajsekharde/ticketing-microservices/proto/payment"
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