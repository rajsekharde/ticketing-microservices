package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
)

type createVenueRequest struct {
	Name string `json:"name"`
	Address string `json:"address"`
}
func createVenue(c *gin.Context) {
	var req createVenueRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	grpcReq := eventpb.CreateVenueRequest{
		Name: req.Name,
		Address: req.Address,
	}

	_, err = event.createVenue(c.Request.Context(), &grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create venue",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "venue created"})
}