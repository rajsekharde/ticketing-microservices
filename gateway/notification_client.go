package main

import (
	"fmt"

	notifpb "github.com/rajsekharde/ticketing-microservices/proto/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type notificationClient struct {
	targetAddr string
	conn *grpc.ClientConn
	stub notifpb.NotificationServiceClient
}

func newNotificationClient(targetAddr string) (*notificationClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Notification service connection failed: %w", err)
	}
	return &notificationClient{
		targetAddr: targetAddr,
		conn: conn,
		stub: notifpb.NewNotificationServiceClient(conn),
	}, nil
}

func (c *notificationClient) close() error {
	return c.conn.Close()
}