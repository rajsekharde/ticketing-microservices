package clients

import (
	"fmt"

	notifpb "github.com/rajsekharde/ticketing-microservices/proto/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationClient struct {
	targetAddr string
	conn *grpc.ClientConn
	Stub notifpb.NotificationServiceClient
}

func NewNotificationClient(targetAddr string) (*NotificationClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Notification service connection failed: %w", err)
	}
	return &NotificationClient{
		targetAddr: targetAddr,
		conn: conn,
		Stub: notifpb.NewNotificationServiceClient(conn),
	}, nil
}

func (c *NotificationClient) Close() error {
	return c.conn.Close()
}