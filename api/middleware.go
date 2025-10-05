package api

import (
	"context"
	"net/http"
	"strings"
)

// Middleware adds authentication and other middleware functionality
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Simple token validation (in production, use proper JWT validation)
		if !isValidToken(token) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract user ID from token (simplified)
		userID := extractUserIDFromToken(token)

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// isValidToken validates a JWT token (simplified implementation)
func isValidToken(token string) bool {
	// In production, use a proper JWT library to validate the token
	// For now, just check if it starts with "jwt_token_"
	return strings.HasPrefix(token, "jwt_token_")
}

// extractUserIDFromToken extracts user ID from token (simplified implementation)
func extractUserIDFromToken(token string) string {
	// In production, decode the JWT token properly
	// For now, extract from the simplified token format
	parts := strings.Split(token, "_")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}
