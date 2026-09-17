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

func newDatabase(cfg *config) (*database, error) {
	// Connection string format: postgres://username:password@host:port/dbname?sslmode=disable
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/user_db?sslmode=disable",
	cfg.db_user, cfg.db_password, cfg.db_host, cfg.db_port)

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

func (d *database) close() {
	d.db.Close()
}

func (d *database) initTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email VARCHAR(255) UNIQUE NOT NULL,
	name VARCHAR(100) NOT NULL,
	role VARCHAR(10) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := d.db.Exec(query)
	return err
}

func (d *database) createUser(user *createUserRequest) error {
	query := `
	INSERT INTO users (email, name, role) VALUES
	($1, $2, $3);
	`

	_, err := d.db.Exec(query, user.Email, user.Name, user.Role)
	return err
}