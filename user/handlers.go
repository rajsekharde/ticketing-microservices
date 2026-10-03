package main

import (
	"context"
	"database/sql"
	"log"

	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"
	// "github.com/rajsekharde/ticketing-microservices/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	user, err := db.getUserQuery(int(req.Id))
	if err != nil {
		log.Printf("[FAILED] Get User: id = %v, error: %v\n", req.Id, err.Error())
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found in database")
		}
		return nil, status.Errorf(codes.Internal, "could not fetch user")
	}

	resUser := userpb.User{
		Id: int64(user.id),
		Email: user.email,
		Name: user.name,
	}
	switch user.role {
	case "ADMIN":
		resUser.Role = userpb.UserRole_ADMIN
	case "CUSTOMER":
		resUser.Role = userpb.UserRole_CUSTOMER
	default:
		resUser.Role = userpb.UserRole_UNSPECIFIED
	}

	log.Printf("Get User: id = %v\n", req.GetId())
	return &resUser, nil
}

func (s *server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
    hash, err := hashPassword(req.Password)
	if err != nil {
        log.Printf("[FAILED] Create User: email = %v, error: %v\n", req.Email, err.Error())
        return nil, status.Errorf(codes.Internal, "failed to hash password")
    }
	
	err = db.createUserQuery(&dbUser{
        email: req.Email,
		password_hash: hash,
        name:  req.Name,
        role:  req.Role.String(),
    })
    if err != nil {
        log.Printf("[FAILED] Create User: email = %v, error: %v\n", req.Email, err.Error())
        // Return a proper gRPC status error (e.g., AlreadyExists if email is taken)
        return nil, status.Errorf(codes.Internal, "failed to insert user into database")
    }

    log.Printf("Create User: email = %v\n", req.Email)
    return &userpb.CreateUserResponse{}, nil // No error field needed in response
}

func (s *server) UserLogin(ctx context.Context, req *userpb.UserLoginRequest) (*userpb.UserLoginResponse, error) {
	hash, err := db.getPasswordQuery(req.Email)
	if err != nil {
		log.Printf("[FAILED] Login User: email = %v, error: %v\n", req.Email, err.Error())
		return nil, status.Errorf(codes.Internal, "failed to fetch password hash")
	}

	if checkPasswordHash(req.Password, hash) == false {
		log.Printf("[FAILED] Login User: email = %v, error: Invalid password\n", req.Email)
		return nil, status.Errorf(codes.Unauthenticated, "password does not match stored hash")
	}

	return &userpb.UserLoginResponse{
		Jwt: "",
	}, nil
}