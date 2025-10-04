package models

// This is io.go (input/output) for json queries and responses

// User details to save to user in DB, received in API call
type UserDetails struct {
	Projects []string `json: "projects"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
}

// HealthCheckResponse returns the health check response duh
type HealthCheckResponse struct {
	Alive bool `json:"alive"`
}

// UserResponse is a general response structure with a status, message and optional json data
type UserResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// ErrorMessageResponse returns the error message response struct
type ErrorMessageResponse struct {
	Response MessageError `json:"response"`
}

// MessageError contains the inner details for the error message response
type MessageError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
