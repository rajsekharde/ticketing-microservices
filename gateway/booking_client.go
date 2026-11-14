package main

import (
	"fmt"

	bookingpb "github.com/rajsekharde/ticketing-microservices/proto/booking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type bookingClient struct {
	targetAddr string
	conn *grpc.ClientConn
	stub bookingpb.BookingServiceClient
}

func newBookingClient(targetAddr string) (*bookingClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Booking service connection failed: %w", err)
	}
	return &bookingClient{
		targetAddr: targetAddr,
		conn: conn,
		stub: bookingpb.NewBookingServiceClient(conn),
	}, nil
}

func (c *bookingClient) close() error {
	return c.conn.Close()
}