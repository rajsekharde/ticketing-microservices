# API Gateway service

## API endpoints
- GET /health : Returns {"message": "API Gateway running"}
- POST /user/register : Creates a new user
- GET /user/:id : Returns user data by ID
- POST /user/login : User login

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