package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Email     string              `json:"email" bson:"email" validate:"required,email"`
	Password  string              `json:"-" bson:"password" validate:"required,min=6"`
	Name      string              `json:"name" bson:"name" validate:"required"`
	Username  string              `json:"username" bson:"username" validate:"required"`
	CompanyID *primitive.ObjectID `json:"companyId,omitempty" bson:"companyId,omitempty"` // Optional company/project link
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt" bson:"updatedAt"`
}

// UserResponse represents the user data sent to the client (without password)
type UserResponse struct {
	ID        primitive.ObjectID  `json:"id"`
	Email     string              `json:"email"`
	Name      string              `json:"name"`
	Username  string              `json:"username"`
	CompanyID *primitive.ObjectID `json:"companyId,omitempty"` // Optional company/project link
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SignupRequest represents the signup request payload
type SignupRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token          string       `json:"token"`
	User           UserResponse `json:"user"`
	CompanyProfile *Company     `json:"companyProfile,omitempty"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Alive bool `json:"alive"`
}

// ErrorMessageResponse represents an error response
type ErrorMessageResponse struct {
	Response MessageError `json:"response"`
}

// MessageError represents an error message
type MessageError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
