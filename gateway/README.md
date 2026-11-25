# API Gateway service

## API endpoints
- GET /health : Returns running status of all services
- POST /users/register : Creates a new customer user. Body: {name, email, password}
- POST /users/login : User login. Body: {email, password}
- POST /users/:id/promote : Promotes a user from role CUSTOMER to ADMIN
- GET /users/:id : Returns user data by ID

- POST /venues : Create venue with seats. Body: {name, address, capacity}
- POST /events : Create event with seats. Body: {name, venue_id, start_time, duration_minutes, ticket_price}
    - start_time: RFC 3339 format : "2026-11-12T19:30:00+05:30"
- GET /events : Returns details of all events
- GET /events/:id/seats : Returns details of all seats of an event

- POST /bookings : Books an event seat. Body: {event_seat_id, event_id, user_id}

Generate JWT token using the user login endpoint, and send the token as Authorization header for protected endpoints. Some endpoints work only for users with role = ADMIN.

## Request flow
Create User:
```bash
Client sends POST /users/register request to Gateway with {name, email, password} as JSON payload
Gateway unmarshalls the JSON body and creates a grpc request for CreateUser using it
CreateUser rpc in user service is invoked with the request struct
user service accepts the request, performs an SQL query in database for user creation, stores password hash and returns the new user id
errors are handled and gRPC status codes are returned by user service
Gateway processes the gRPC response, and sends an HTTP response to client accordingly
```

User Login
```bash
Client sends POST /users/login request to Gateway with {email, password} as JSON payload
Gateway invokes the rpc for user login using the credentials as request
user service checks password with stored hash
A JWT is generated using the user id and role, signed with a secret key
The JWT is sent as grpc response to Gateway, which is then forwarded to the client with a HTTP response
```

Get user details by ID
```bash
Client logs in by generating a JWT
Client sends GET /users/:id request to Gateway with the JWT as Authorization header
Gateway parses the JWT and verifies the token's validity
Caller user's id and role is extracted from the JWT claims. Regular users can only check their own details, and admins can check other users' details
Gateway sends grpc request to user service
user service performs database query and returns the user's data
Gateway sends the user data as HTTP response
```

## GET /health endpoint

All the services expose a health/ping endpoint. The gateway sends request to all these endpoints, and returns all the running statuses with timestamps. Results for a certain interval are cached, like 10 seconds so repeated health checks don't overwhelm the application.