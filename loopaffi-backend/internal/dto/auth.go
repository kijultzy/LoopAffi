package dto

// LoginRequest is the payload for POST /api/v1/auth/login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest is the payload for POST /api/v1/auth/register
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Whatsapp string `json:"whatsapp" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse is the success response for login
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type ForgotPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UserProfile is a safe user representation (no password)
type UserProfile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
