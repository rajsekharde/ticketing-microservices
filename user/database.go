package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	// userpb "github.com/rajsekharde/ticketing-microservices/proto/user"
	// "github.com/rajsekharde/ticketing-microservices/shared"
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
	password_hash VARCHAR(255) NOT NULL,
	name VARCHAR(100) NOT NULL,
	role VARCHAR(20) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := d.db.Exec(query)
	return err
}

type dbUser struct {
	email string
	password_hash string
	name string
	role string
}

func (d *database) createUserQuery(user *dbUser) error {
	query := `
	INSERT INTO users (email, password_hash, name, role) VALUES
	($1, $2, $3, $4);
	`

	_, err := d.db.Exec(query, user.email, user.password_hash, user.name, user.role)
	return err
}

type user struct {
	id int
	email string
	name string
	role string
}
func (d *database) getUserQuery(id int) (*user, error) {
	query := `
	SELECT id, email, name, role FROM users WHERE id = $1;
	`

	var user user
	err := d.db.QueryRow(query, id).Scan(&user.id, &user.email, &user.name, &user.role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Returns the stored password hash for user with given email
func (d *database) getPasswordQuery(email string) (string, error) {
	query := `
	SELECT password_hash FROM users WHERE email = $1;
	`
	var hash string
	err := d.db.QueryRow(query, email).Scan(&hash)
	if err != nil {
		return "", err
	}
	
	return hash, nil;
}