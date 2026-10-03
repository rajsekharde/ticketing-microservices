package shared

// Data received by gateway for creating a user
type CreateUserRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
}

// Data sent to client for GET /users/:id
type GetUserResponse struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// gateway request for user login
type UserLoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}