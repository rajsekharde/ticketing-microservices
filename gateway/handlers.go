package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Returns status of server
func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
	})
}