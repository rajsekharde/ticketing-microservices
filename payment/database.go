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
	db  *sql.DB
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
		db:  db,
	}, nil
}

// Closes the DB connection
func (d *database) close() {
	d.db.Close()
}

// Creates DB tables for payment service
func (d *database) initTables() error {
	query := `CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		booking_id INT NOT NULL,
		amount INT NOT NULL,
		status VARCHAR(20) DEFAULT 'PENDING',
		completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Create a new transaction
func (d *database) newTransactionQuery(bookingId int, amount int) (int, error) {
	query := `
	INSERT INTO transactions (booking_id, amount)
	VALUES ($1, $2)
	RETURNING id;
	`

	var transId int
	err := d.db.QueryRow(query, bookingId, amount).Scan(&transId)
	if err != nil {
		return transId, err
	}

	return transId, nil
}

// Update transaction status
func (d *database) updateTransStatusQuery(transId int, newStatus string) error {
	query := `
	UPDATE transactions
	SET status = $1
	WHERE id = $2;
	`

	_, err := d.db.Exec(query, newStatus, transId)
	if err != nil {
		return err
	}

	return nil
}