package clients

import (
	"fmt"

	paymentpb "github.com/rajsekharde/ticketing-microservices/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	targetAddr string
	conn *grpc.ClientConn
	Stub paymentpb.PaymentServiceClient
}

func NewPaymentClient(targetAddr string) (*PaymentClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Payment service connection failed: %w", err)
	}
	return &PaymentClient{
		targetAddr: targetAddr,
		conn: conn,
		Stub: paymentpb.NewPaymentServiceClient(conn),
	}, nil
}

func (c *PaymentClient) Close() error {
	return c.conn.Close()
}