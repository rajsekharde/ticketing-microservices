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

func (c *EventClient) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Stub.CreateEvent(ctx, req)
}

func (c *EventClient) GetEvents(ctx context.Context, req *eventpb.GetEventsRequest) (*eventpb.GetEventsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Stub.GetEvents(ctx, req)
}

func (c *EventClient) GetEventSeats(ctx context.Context, req *eventpb.GetEventSeatsRequest) (*eventpb.GetEventSeatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Stub.GetEventSeats(ctx, req)
}

func (c *EventClient) GetVenues(ctx context.Context, req *eventpb.GetVenuesRequest) (*eventpb.GetVenuesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Stub.GetVenues(ctx, req)
}