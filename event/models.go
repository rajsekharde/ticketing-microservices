package main

import "time"

type venue struct {
	name string
	address string
	capacity int
}

type event struct {
	id int
	name string
	venueId int
	startTime time.Time
	durationMinutes int
	ticketPrice int

	venueName string
	venueAddress string
}