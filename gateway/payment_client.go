package main

import (
	"fmt"

	paymentpb "github.com/rajsekharde/ticketing-microservices/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type paymentClient struct {
	targetAddr string
	conn *grpc.ClientConn
	stub paymentpb.PaymentServiceClient
}

func newPaymentClient(targetAddr string) (*paymentClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Payment service connection failed: %w", err)
	}
	return &paymentClient{
		targetAddr: targetAddr,
		conn: conn,
		stub: paymentpb.NewPaymentServiceClient(conn),
	}, nil
}

func (c *paymentClient) close() error {
	return c.conn.Close()
}