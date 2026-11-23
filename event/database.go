package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type database struct {
	url string
	db *sql.DB
}

// Creates a new database connection. Uses values from config struct to create DB URL
func newDatabase(cfg *config) (*database, error) {
	// Connection string format: postgres://username:password@host:port/dbname?sslmode=disable
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	cfg.db_user, cfg.db_password, cfg.db_host, cfg.db_port, cfg.db_name)

	// intialize datatbase handle, registers configuration, driver, and connection string
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// tests the network connection and credentials
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	return &database{
		url: url,
		db: db,
	}, nil
}

// Closes the DB connection
func (d *database) close() {
	d.db.Close()
}

// Creates DB tables for event service
func (d *database) initTables() error {
	queries := []string {
		`CREATE TABLE IF NOT EXISTS venues (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		address VARCHAR(100) UNIQUE NOT NULL,
		capacity INT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS seats (
		id SERIAL PRIMARY KEY,
		venue_id INT NOT NULL,
		seat_number VARCHAR(10) NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		venue_id INT NOT NULL,
		name VARCHAR(100) NOT NULL,
		start_time TIMESTAMP NOT NULL,
		duration_minutes INT NOT NULL,
		ticket_price INT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS event_seats (
		id SERIAL PRIMARY KEY,
		event_id INT NOT NULL,
		seat_id INT NOT NULL,
		status VARCHAR(20) DEFAULT 'AVAILABLE',
		ticket_price INT NOT NULL
		)`,
	}

	for _, query := range queries {
		_, err := d.db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// Adds a new venue to venues table.
// Also adds seats for the venue in seats table with serial seat numbers
func (d *database) createVenueWithSeatsQuery(v *venue) (int, error) {
	venueQuery := `
	INSERT INTO venues (name, address, capacity) VALUES ($1, $2, $3) RETURNING id;
	`
	seatQuery := `
	INSERT INTO seats (venue_id, seat_number) VALUES ($1, $2);
	`

	var venueId int

	// start a transaction
	tx, err := d.db.Begin()
	if err != nil {
		return venueId, err
	}
	// rollback the transaction if any error occurs
	defer tx.Rollback()

	err = tx.QueryRow(venueQuery, v.name, v.address, v.capacity).Scan(&venueId)
	if err != nil {
		return venueId, err
	}

	for i := range v.capacity {
		_, err := tx.Exec(seatQuery, venueId, strconv.Itoa(i+1))
		if err != nil {
			return venueId, err
		}
	}

	// commit the transaction
	return venueId, tx.Commit()
}

// Creates a new event in events table.
// Also maps the venue seats with the event and creates entries in event_seats table.
func (d *database) createEventWithSeatsQuery(e *event) (int, error) {
	var eventId int

	venueQuery := `
	SELECT id FROM seats WHERE venue_id = $1;
	`
	eventQuery := `
	INSERT INTO events
	(name, venue_id, start_time, duration_minutes, ticket_price)
	VALUES ($1, $2, $3, $4, $5) returning id;
	`
	eventSeatQuery := `
	INSERT INTO event_seats (event_id, seat_id, ticket_price) VALUES ($1, $2, $3);
	`

	// get all seats for the venue_id
	rows, err := d.db.Query(venueQuery, e.venueId)
	if err != nil {
		return eventId, err
	}
	defer rows.Close()

	// store all seat IDs in an array
	var seatIDs []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return eventId, err
		}
		seatIDs = append(seatIDs, id)
	}
	// check for invalid venue / no seats
	if len(seatIDs) == 0 {
		return eventId, fmt.Errorf("venue_id %d has no seats configured", e.venueId)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return eventId, err
	}
	defer tx.Rollback()

	// create the event
	err = tx.QueryRow(eventQuery, e.name, e.venueId, e.startTime, e.durationMinutes, e.ticketPrice).Scan(&eventId)
	if err != nil {
		return eventId, err
	}
	
	// create the event seats
	for _, id := range seatIDs {
		_, err = tx.Exec(eventSeatQuery, eventId, id, e.ticketPrice)
		if err != nil {
			return eventId, err
		}
	}

	return eventId, tx.Commit()
}

// Returns details of all events
func (d *database) getEventsQuery() ([]event, error) {
	var events []event

	eventsQuery := `
	SELECT
	e.id, e.name, v.name, v.address, e.start_time, e.duration_minutes, e.ticket_price
	FROM events AS e
	INNER JOIN venues AS v
	ON e.venue_id = v.id;
	`

	rows, err := d.db.Query(eventsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e event
		err := rows.Scan(
			&e.id,
			&e.name,
			&e.venueName,
			&e.venueAddress,
			&e.startTime,
			&e.durationMinutes,
			&e.ticketPrice)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

// Returns all seats of an event
func (d *database) getEventSeatsQuery(eventId int) ([]eventSeat, error) {
	query := `
	SELECT id, event_id, seat_id, status, ticket_price
	FROM event_seats WHERE event_id = $1;
	`

	var eventSeats []eventSeat

	rows, err := d.db.Query(query, eventId)
	if err != nil {
		return eventSeats, err
	}
	defer rows.Close()

	for rows.Next() {
		var e eventSeat
		err = rows.Scan(
			&e.id,
			&e.eventId,
			&e.seatId,
			&e.status,
			&e.price,
		)
		if err != nil {
			return eventSeats, err
		}
		eventSeats = append(eventSeats, e)
	}

	return eventSeats, nil
}

// Returns details of all venues
func (d *database) getVenuesQuery() ([]venue, error) {
	var venues []venue

	venueQuery := `SELECT id, name, address, capacity FROM venues;`

	rows, err := d.db.Query(venueQuery)
	if err != nil {
		return venues, err
	}
	defer rows.Close()

	for rows.Next() {
		var v venue
		err = rows.Scan(
			&v.id,
			&v.name,
			&v.address,
			&v.capacity,
		)
		if err != nil {
			return venues, err
		}
		venues = append(venues, v)
	}

	return venues, nil
}