package handlers

import (
	"github.com/rajsekharde/ticketing-microservices/gateway/clients"
)

type Clients struct {
	User *clients.UserClient
	Event *clients.EventClient
	Booking *clients.BookingClient
	Payment *clients.PaymentClient
	Notification *clients.NotificationClient
}

var cl Clients

func LoadClients(c *Clients) {
	cl.User = c.User
	cl.Event = c.Event
	cl.Booking = c.Booking
	cl.Payment = c.Payment
	cl.Notification = c.Notification
}