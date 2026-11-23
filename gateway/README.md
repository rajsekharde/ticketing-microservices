# API Gateway service

## API endpoints
- GET /health : Returns running status of all services
- POST /users/register : Creates a new user. Body: {name, email, password, role}
    - Roles: ADMIN, CUSTOMER
- GET /users/:id : Returns user data by ID
- POST /users/login : User login. Body: {email, password}
- POST /venues : Create venue with seats. Body: {name, address, capacity}
- POST /events : Create event with seats. Body: {name, venue_id, start_time, duration_minutes, ticket_price}
    - start_time: RFC 3339 format : "2026-11-12T19:30:00+05:30"
- GET /events : Returns details of all events
- GET /events/:id/seats : Returns details of all seats of an event

## Request flow
Create User:
```bash
Client sends POST /user to Gateway with credentials as JSON payload
Gateway unmarshalls the JSON body into a CreateUserRequest struct
A userpb.CreateRequestUser struct is created using the credentials, as req
CreateUser rpc in user service is invoked with the req struct
user service accepts the request, and performs an SQL query in database for user creation
gRPC status codes are returned by user service, according to errors
Gateway processes the gRPC response, and sends an HTTP response to client accordingly
```

## GET /health endpoint

All the services expose a health/ping endpoint. The gateway sends request to all these endpoints, and returns all the running statuses with timestamps. Results for a certain interval are cached, like 10 seconds so repeated health checks don't overwhelm the application.