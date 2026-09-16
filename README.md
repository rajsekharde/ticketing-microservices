# Go microservices backend for event ticket booking

## Tech stack
- Go for core service logic
- gRPC for inter-service communication
- RabbitMQ as message broker
- Docker for containerization
- kind for local Kubernetes deployment

## Services

### User
Handles:
- Registration
- Login
- Password hashing
- Access, refresh tokens
- User roles

### Event
Handles:
- Data of event, venue, shows
- Seat layouts

### Booking
Manages the lifecycle of a ticket reservation

Handles:
- bookings
- seat reservations
- reservation expiration
- booking status

### Payment
Handles payment attempts.
Accepts payment request with credentials. Returns payment status, metadata as response.

### Notification
Sends notification to users when an event occurs, like booking confirmed, cancelled or failed.
Should be asynchronous to booking lifecycle.


## Supporting infrastructure
- API Gateway: Terminates HTTP, validates JWT, rate limiting, routes to internal gRPC services
- Redis: Caching
- RabbitMQ: Async events
- PostgreSQL: Stores structured data of services in separate databases


## Complete booking flow

1. Registration
```bash
user sends credentials to api gateway
api forwards the request to user service
user service creates a new user in db
registration status is send to user via response from api
```

2. Login
```bash
user sends credentials to api
api forwards the request to user service
user service verifies the credentials, creates access & refresh tokens
access token is sent to user via response from api
user is logged in, and event data / dashboard are displayed, based on the role of user
```

3. Browse events
```bash
event data is fetched from event service
```

4. Book ticket
```bash
event date, venue, time, seats are selected and sent to booking service as request via api
booking service performs a transaction:
- verify avilability of seats
- temporarily lock seats for a duration
- send request to payment service with credentials
```

5. Payment
```bash
payment service verifies credentials, completes payment
sends payment status, metadata to booking service
```

6. Booking confirmation
```bash
booking service books / frees seats based on payment status
sends response to user via api
publishes event to message queue
```

7. Notification
```bash
notification service reads event from message queue
sends event notification to user, async to booking lifecycle
```

Possible failures:
- Registration / Login fails due to wrong credentials, request timeout / network issues, server errors. Handling: user is logged out
- More than one user tries to book a seat at the same time. Handled by database transactions, atomic operations
- Payment failure due to wrong credentials, timeout. Handled by releasing temporarily booked seats
- Payment done more than once: each payment should be unique and traceable
- Notification failure. Handled by retries, booking status unaffected by notifications


## Project structure
```bash
ticketing-microservices/
    gateway/
        main.go
        handlers.go
        models.go
        ...
    user/
        main.go
        database.go
        models.go
        ...
    event/
    booking/
    payment/
    notification/
    bin/
    
    docker/
        gateway.Dockerfile
        user.Dockerfile
        ...
    k8s/
    proto/
    
    go.mod
    go.sum
    docker-compose.yml
    Makefile
    README.md
```