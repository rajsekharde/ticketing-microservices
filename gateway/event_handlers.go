package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type createVenueRequest struct {
	Name string `json:"name"`
	Address string `json:"address"`
	Capacity int `json:"capacity"`
}
func createVenue(c *gin.Context) {
	var req createVenueRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	grpcReq := eventpb.CreateVenueRequest{
		Name: req.Name,
		Address: req.Address,
		Capacity: int64(req.Capacity),
	}

	grpcResp, err := event.createVenue(c.Request.Context(), &grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create venue",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "venue created",
		"venue_id": grpcResp.VenueId,
	})
}

type createEventRequest struct {
	Name string `json:"name"`
	VenueId int `json:"venue_id"`
	StartTime time.Time `json:"start_time"`
	DurationMinutes int `json:"duration_minutes"`
	TicketPrice int `json:"ticket_price"`
}
func createEvent(c *gin.Context) {
	var req createEventRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	grpcReq := eventpb.CreateEventRequest{
		Name: req.Name,
		VenueId: int64(req.VenueId),
		StartTime: timestamppb.New(req.StartTime),
		DurationMinutes: int64(req.DurationMinutes),
		TicketPrice: int64(req.TicketPrice),
	}

	grpcResp, err := event.createEvent(c.Request.Context(), &grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create event",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "event created",
		"event_id": grpcResp.EventId,
	})
}