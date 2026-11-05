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

// Creates DB tables for event service
func (d *database) initTables() error {
	queries := []string {
		`CREATE TABLE IF NOT EXISTS venues (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		address VARCHAR(100) NOT NULL
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
		time TIMESTAMP NOT NULL,
		duration INTERVAL
		)`,
		`CREATE TABLE IF NOT EXISTS event_seats (
		id SERIAL PRIMARY KEY,
		event_id INT NOT NULL,
		seat_id INT NOT NULL,
		status VARCHAR(20) DEFAULT 'AVAILABLE',
		price INT NOT NULL
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