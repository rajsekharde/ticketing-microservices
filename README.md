# A Go microservices backend for event ticket booking

## Tech stack
- Go for core service logic
- gRPC for inter-service communication
- RabbitMQ as message broker
- Docker for containerization
- kind for local Kubernetes deployment

## Services
- API Gateway: Terminates HTTP, validates JWT, rate limiting, routes to internal gRPC services
- Auth: Registration, login, JWT issuance/refresh, logout/revocation
- User: Operations related to user profile data
- Event: Operations related to event, venue, shows etc
- Booking: Handles seat reservation requests, creates booking record, etc
- Payment: Payment processing
- Notification: Sends confirmation/cancellation emails

## Supporting infrastructure
- Redis: Caching, seat locks
- RabbitMQ: Async events
- PostgreSQL: Stores structured data of services in separate databases