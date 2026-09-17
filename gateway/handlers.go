package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Returns status of server
func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
	})
}

// Returns user data
func getUserById(c *gin.Context) {
	// 1. Extract path parameter ":id"
	idParam := c.Param("id")

	// 2. Parse string ID to int64
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID, must be a positive integer",
		})
		return
	}

	// 3. Call the gRPC client wrapper (pass c.Request.Context() for cancellation propagation)
	user, err := user.getUser(c.Request.Context(), userID)
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
	c.JSON(http.StatusOK, gin.H{
		"id":    user.GetId(),
		"email": user.GetEmail(),
		"name":  user.GetName(),
		"role":  user.GetRole().String(), // Converts enum to readable string (e.g. "USER_ROLE_CUSTOMER")
	})
}
