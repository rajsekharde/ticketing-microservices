package clients

import (
	"context"
	"fmt"
	"time"

	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	targetAddr string
	conn *grpc.ClientConn
	Stub userpb.UserServiceClient
}

func NewUserClient(targetAddr string) (*UserClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("User service connection failed: %w", err)
	}
	return &UserClient{
		targetAddr: targetAddr,
		conn: conn,
		Stub: userpb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) GetUser(ctx context.Context, id int64) (*userpb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Stub.GetUser(ctx, &userpb.GetUserRequest{Id: id})
}

func (c *UserClient) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Stub.CreateUser(ctx, req)
}

func (c *UserClient) UserLogin(ctx context.Context, req *userpb.UserLoginRequest) (*userpb.UserLoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Stub.UserLogin(ctx, req)
}