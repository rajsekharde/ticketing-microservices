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

type eventStruct struct {
	Id int `json:"event_id"`
	Name string `json:"event_name"`
	VenueName string `json:"venue_name"`
	VenueAddress string `json:"venue_address"`
	StartTime time.Time `json:"start_time"`
	DurationMinutes int `json:"duration_minutes"`
	TicketPrice int `json:"ticket_price"`
}
func getEvents(c *gin.Context) {
	grpcReq := eventpb.GetEventsRequest{}

	grpcResp, err := event.getEvents(c.Request.Context(), &grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch events",
		})
		return
	}
	if len(grpcResp.Events) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "no events found",
		})
		return
	}

	var events []eventStruct
	for _, e := range grpcResp.Events {
		ev := eventStruct{
			Id: int(e.Id),
			Name: e.Name,
			VenueName: e.VenueName,
			VenueAddress: e.VenueAddress,
			StartTime: e.StartTime.AsTime(),
			DurationMinutes: int(e.DurationMinutes),
			TicketPrice: int(e.TicketPrice),
		}

		events = append(events, ev)
	}

	c.JSON(http.StatusOK, events)
}