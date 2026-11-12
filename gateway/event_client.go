package main

import (
	"context"
	"fmt"
	"time"

	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type eventClient struct {
	targetAddr string
	conn *grpc.ClientConn
	stub eventpb.EventServiceClient
}

func newEventClient(targetAddr string) (*eventClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Event service connection failed: %w", err)
	}
	return &eventClient{
		targetAddr: targetAddr,
		conn: conn,
		stub: eventpb.NewEventServiceClient(conn),
	}, nil
}

func (c *eventClient) close() error {
	return c.conn.Close()
}

func (c *eventClient) createVenue(ctx context.Context, req *eventpb.CreateVenueRequest) (*eventpb.CreateVenueResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.stub.CreateVenue(ctx, req)
}

func (c *eventClient) createEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.stub.CreateEvent(ctx, req)
}

func (c *eventClient) getEvents(ctx context.Context, req *eventpb.GetEventsRequest) (*eventpb.GetEventsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.stub.GetEvents(ctx, req)
}