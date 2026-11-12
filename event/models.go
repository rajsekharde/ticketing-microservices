package main

import "time"

type venue struct {
	name string
	address string
	capacity int
}

type event struct {
	name string
	venueId int
	startTime time.Time
	durationMinutes int
}