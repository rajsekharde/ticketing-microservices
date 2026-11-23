package clients

import (
	"fmt"

	bookingpb "github.com/rajsekharde/ticketing-microservices/proto/booking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookingClient struct {
	targetAddr string
	conn *grpc.ClientConn
	Stub bookingpb.BookingServiceClient
}

func NewBookingClient(targetAddr string) (*BookingClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Booking service connection failed: %w", err)
	}
	return &BookingClient{
		targetAddr: targetAddr,
		conn: conn,
		Stub: bookingpb.NewBookingServiceClient(conn),
	}, nil
}

func (c *BookingClient) Close() error {
	return c.conn.Close()
}