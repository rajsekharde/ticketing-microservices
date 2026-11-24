package main

import (
	"context"
	"database/sql"
	"fmt"
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

// Creates DB tables for booking service
func (d *database) initTables() error {
	queries := []string {
		`CREATE TABLE IF NOT EXISTS bookings (
		id SERIAL PRIMARY KEY,
		event_id INT NOT NULL,
		user_id INT NOT NULL,
		status VARCHAR(20) DEFAULT 'PENDING',
		booked_at TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS booking_items (
		id SERIAL PRIMARY KEY,
		booking_id INT NOT NULL,
		event_seat_id INT NOT NULL
		);`,
	}

	for _, query := range queries {
		_, err := d.db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// Creates a new booking entry along with its booking items
func (d *database) createBookingQuery(eventSeatIds []int, eventId int, userId int) (int, error) {
	bookingQuery := `
	INSERT INTO bookings (event_id, user_id)
	VALUES ($1, $2)
	RETURNING id;
	`
	itemQuery := `
	INSERT INTO booking_items (booking_id, event_seat_id)
	VALUES ($1, $2);
	`

	var bookingId int

	tx, err := d.db.Begin()
	if err != nil {
		return bookingId, err
	}
	defer tx.Rollback()

	// create booking
	err = tx.QueryRow(bookingQuery, eventId, userId).Scan(&bookingId)
	if err != nil {
		return bookingId, err
	}
	
	// create booking items
	for _, id := range eventSeatIds {
		_, err = tx.Exec(itemQuery, bookingId, id)
		if err != nil {
			return bookingId, err
		}
	}

	return bookingId, tx.Commit()
}

// Updates booking status
func (d *database) updateBookingStatusQuery(bookingId int, newStatus string) error {
	query := `
	UPDATE bookings
	SET status = $1
	WHERE id = $2 AND status = 'PENDING';
	`

	_, err := d.db.Exec(query, newStatus, bookingId)
	if err != nil {
		return err
	}

	return nil
}