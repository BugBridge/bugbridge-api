package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/BugBridge/bugbridge-api/databases"
)

// Middleware adds database context and JWT authentication
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add database to context
		dbHelper := r.Context().Value("dbHelper").(databases.DatabaseHelper)
		ctx := context.WithValue(r.Context(), "dbHelper", dbHelper)

		// Extract JWT token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-super-secret-jwt-key-change-this-in-production"), nil
		})

		if err != nil {
			zap.S().With(err).Error("Token parsing failed")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx = context.WithValue(ctx, "userID", claims.UserID)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DatabaseMiddleware adds database to request context
func DatabaseMiddleware(dbHelper databases.DatabaseHelper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "dbHelper", dbHelper)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
