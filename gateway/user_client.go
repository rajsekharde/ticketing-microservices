package main

import (
	"context"
	"fmt"
	"time"

	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userClient struct {
	targetAddr string
	conn *grpc.ClientConn
	stub userpb.UserServiceClient
}

func newUserClient(targetAddr string) (*userClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("User service connection failed: %w", err)
	}
	return &userClient{
		targetAddr: targetAddr,
		conn: conn,
		stub: userpb.NewUserServiceClient(conn),
	}, nil
}

func (c *userClient) close() error {
	return c.conn.Close()
}

func (c *userClient) getUser(ctx context.Context, id int64) (*userpb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.stub.GetUser(ctx, &userpb.GetUserRequest{Id: id})
}

func (c *userClient) createUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.stub.CreateUser(ctx, req)
}