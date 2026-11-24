package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rajsekharde/ticketing-microservices/booking/clients"
)

type config struct {
	port string
	db_host string
	db_port string
	db_user string
	db_password string
	db_name string
	eventAddress string
	paymentAddress string
}
func loadEnv(path string) config {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("Failed to load env file. Using environment variables")
	}

	port := os.Getenv("BOOKING_PORT")
	if port == "" {
		port = "50053"
	}

	eventHost := os.Getenv("EVENT_HOST")
	eventPort := os.Getenv("EVENT_PORT")
	eventAddress := eventHost + ":" + eventPort

	paymentHost := os.Getenv("PAYMENT_HOST")
	paymentPort := os.Getenv("PAYMENT_PORT")
	paymentAddress := paymentHost + ":" + paymentPort

	return config{
		port: port,
		db_host: os.Getenv("POSTGRES_HOST"),
		db_port: os.Getenv("POSTGRES_PORT"),
		db_user: os.Getenv("POSTGRES_USER"),
		db_password: os.Getenv("POSTGRES_PASSWORD"),
		db_name: os.Getenv("BOOKING_DB"),
		eventAddress: eventAddress,
		paymentAddress: paymentAddress,
	}
}

type grpcClients struct {
	event clients.EventClient
	payment clients.PaymentClient
}

func createClients(cgf *config) (*grpcClients, error) {
	event, err := clients.NewEventClient(cfg.eventAddress)
	if err != nil {
		return nil, err
	}

	payment, err := clients.NewPaymentClient(cfg.paymentAddress)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		event: *event,
		payment: *payment,
	}, nil
}