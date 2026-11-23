package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	eventpb "github.com/rajsekharde/ticketing-microservices/proto/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateVenueRequest struct {
	Name string `json:"name"`
	Address string `json:"address"`
	Capacity int `json:"capacity"`
}
func CreateVenue(c *gin.Context) {
	var req CreateVenueRequest

	// check user role from claims
	valRole, exists := c.Get("user_role")
	userRole, ok := valRole.(string)
	if !exists || !ok || userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if userRole != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

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

	grpcResp, err := cl.Event.CreateVenue(c.Request.Context(), &grpcReq)
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

type CreateEventRequest struct {
	Name string `json:"name"`
	VenueId int `json:"venue_id"`
	StartTime time.Time `json:"start_time"`
	DurationMinutes int `json:"duration_minutes"`
	TicketPrice int `json:"ticket_price"`
}
func CreateEvent(c *gin.Context) {
	var req CreateEventRequest

	// check user role from claims
	valRole, exists := c.Get("user_role")
	userRole, ok := valRole.(string)
	if !exists || !ok || userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if userRole != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

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

	grpcResp, err := cl.Event.CreateEvent(c.Request.Context(), &grpcReq)
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

type EventStruct struct {
	Id int `json:"event_id"`
	Name string `json:"event_name"`
	VenueName string `json:"venue_name"`
	VenueAddress string `json:"venue_address"`
	StartTime time.Time `json:"start_time"`
	DurationMinutes int `json:"duration_minutes"`
	TicketPrice int `json:"ticket_price"`
}
func GetEvents(c *gin.Context) {
	grpcReq := eventpb.GetEventsRequest{}

	grpcResp, err := cl.Event.GetEvents(c.Request.Context(), &grpcReq)
	if err != nil || grpcResp == nil {
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

	var events []EventStruct
	for _, e := range grpcResp.Events {
		ev := EventStruct{
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

type EventSeatStruct struct {
	Id int `json:"event_seat_id"`
	EventId int `json:"event_id"`
	SeatId int `json:"seat_id"`
	Status string `json:"status"`
	Price int `json:"price"`
}
func GetEventSeats(c *gin.Context) {
	eventIdParam := c.Param("id")
	eventId, err := strconv.ParseInt(eventIdParam, 10, 64) // 64-bit decimal
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Invalid event id. Must be a positive integer",
		})
		return
	}

	grpcReq := eventpb.GetEventSeatsRequest{EventId: eventId}

	grpcResp, err := cl.Event.GetEventSeats(c.Request.Context(), &grpcReq)
	if err != nil || grpcResp == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch event seats",
		})
	}
	if len(grpcResp.Seats) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message":"event seats not found",
		})
	}

	var seats []EventSeatStruct
	for _, s := range grpcResp.Seats {
		e := EventSeatStruct{
			Id: int(s.Id),
			EventId: int(s.EventId),
			SeatId: int(s.SeatId),
			Status: s.Status,
			Price: int(s.Price),
		}
		seats = append(seats, e)
	}

	c.JSON(http.StatusOK, seats)
}