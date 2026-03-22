package dto

import "time"

// LoginRequest - untuk login request
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" validate:"required,min=3,max=255"`
	Password        string `json:"password" validate:"required,min=6,max=255"`
}

// RegisterRequest - untuk register request
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=255"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
}

// UserDataResponse - untuk user response
type UserDataResponse struct {
	ID         int64     `json:"id"`
	UniqueID   string    `json:"unique_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	// CreatedAt time.Time `json:"created_at"` // ← tambah ini juga
}

// LoginResponse - untuk login response
type LoginResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresIn    int              `json:"expires_in"`
	User         UserDataResponse `json:"user"`
}

// RegisterResponse - untuk register response
type RegisterResponse struct {
	ID        int64     `json:"id"`
	UniqueID  string    `json:"unique_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ChangePasswordRequest
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6,max=255"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=255"`
}

type AuthResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresIn    int              `json:"expires_in"`
	User         UserDataResponse `json:"user"`
}
