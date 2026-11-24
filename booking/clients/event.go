// gRPC client for event service

package clients

import (
	"context"
	"fmt"
	"time"

	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventClient struct {
	targetAddr string
	conn *grpc.ClientConn
	Stub eventpb.EventServiceClient
}

func NewEventClient(targetAddr string) (*EventClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Event service connection failed: %w", err)
	}
	return &EventClient{
		targetAddr: targetAddr,
		conn: conn,
		Stub: eventpb.NewEventServiceClient(conn),
	}, nil
}

func (c *EventClient) Close() error {
	return c.conn.Close()
}

func (c *EventClient) CreateVenue(ctx context.Context, req *eventpb.CreateVenueRequest) (*eventpb.CreateVenueResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Stub.CreateVenue(ctx, req)
}

func (c *EventClient) LockSeat(ctx context.Context, req *eventpb.LockSeatRequest) (*eventpb.LockSeatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Stub.LockSeat(ctx, req)
}