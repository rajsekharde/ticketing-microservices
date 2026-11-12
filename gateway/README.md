# API Gateway service

## API endpoints
- GET /health : Returns {"message": "API Gateway running"}
- POST /user/register : Creates a new user. Body: {name, email, password, role}
- GET /user/:id : Returns user data by ID
- POST /user/login : User login
- POST /venue : Create venue. Body: {name, address, capacity}
- POST /venue : Create event. Body: {name, venue_id, start_time, duration_minutes}
    - start_time: RFC 3339 format : "2026-11-12T19:30:00+05:30"

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