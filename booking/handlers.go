package main

import (
	"context"

	bookingpb "github.com/rajsekharde/ticketing-microservices/proto/booking"
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