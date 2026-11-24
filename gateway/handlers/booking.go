package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bookingpb "github.com/rajsekharde/ticketing-microservices/proto/booking"
)

// Books an event seat
type BookSeatRequest struct {
	EventSeatId int `json:"event_seat_id"`
	EventId int `json:"event_id"`
	UserId int `json:"user_id"`
}
func BookSeat(c *gin.Context) {
	var req BookSeatRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	grpcReq := bookingpb.BookSeatRequest{
		EventSeatId: int64(req.EventSeatId),
		EventId: int64(req.EventId),
		UserId: int64(req.UserId),
	}

	grpcResp, err := cl.Booking.Stub.BookSeat(c.Request.Context(), &grpcReq)
	if err != nil || grpcResp == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to book seat",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "seat booked",
		"booking_id": grpcResp.BookingId,
	})
}