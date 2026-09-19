package shared

// Data received by gateway for creating a user
type CreateUserRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Role string `json:"role"`
}