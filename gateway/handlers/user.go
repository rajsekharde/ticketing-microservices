package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "github.com/rajsekharde/ticketing-microservices/proto/user"
	"github.com/rajsekharde/ticketing-microservices/shared"
)


// Returns user data
func GetUserById(c *gin.Context) {
	// 1. Extract path parameter ":id"
	idParam := c.Param("id")

	// 2. Parse string ID to int64
	userID, err := strconv.ParseInt(idParam, 10, 64)  // 64-bit decimal
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID, must be a positive integer",
		})
		return
	}

	// 3. Call the gRPC client wrapper (pass c.Request.Context() for cancellation propagation)
	user, err := cl.User.GetUser(c.Request.Context(), userID)
	if err != nil {
		// Map gRPC status codes to appropriate HTTP status codes
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user details"})
		return
	}

	// 4. Return user JSON response
	c.JSON(http.StatusOK, shared.GetUserResponse{
		Id:    int(user.GetId()),
		Email: user.GetEmail(),
		Name:  user.GetName(),
		Role:  user.GetRole().String(), // Converts enum to readable string (e.g. "USER_ROLE_CUSTOMER")
	})
}

// Creates a new user
func CreateUser(c *gin.Context) {
	var req shared.CreateUserRequest

    // Bind request body to the struct and check for errors
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

	grpcReq := userpb.CreateUserRequest{
		Email: req.Email,
		Name: req.Name,
		Password: req.Password,
	}
	switch req.Role {
	case "CUSTOMER":
		grpcReq.Role = userpb.UserRole_CUSTOMER
	case "ADMIN":
		grpcReq.Role = userpb.UserRole_ADMIN
	default:
		grpcReq.Role = userpb.UserRole_UNSPECIFIED
	}

	_, err := cl.User.CreateUser(c.Request.Context(), &grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":    "user created",
	})
}

// Handles user login
func UserLogin(c *gin.Context) {
	var req shared.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	grpcReq := userpb.UserLoginRequest{
		Email: req.Email,
		Password: req.Password,
	}
	resp, err := cl.User.UserLogin(c.Request.Context(), &grpcReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": resp.Jwt,
	})
}