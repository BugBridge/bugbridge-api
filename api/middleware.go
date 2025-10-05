
package api

import (
	"net/http"
	"context"
	"fmt"
	"strings"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Middleware adds some basic header authentication around accessing the routes
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do Authentication stuff?
		// use r.URL to get url
		next.ServeHTTP(w, r)
	})
}	

type ctxKey string

const userIDKey ctxKey = "user_id"

//Check JWT tokens and verify them 
func authMiddleware(next http.Handler) http.Handler
{
	secret := os.Getenv("SECRET")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)
	{
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))

		//if the header is empty Error
		if authHeader == ""
		{
			http.Error(w, "No authorization Header", http.StatusUnauthorized)
			return
		}

		//If there is not 2 parts to header Error
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer")
		{
			http.Error(w, "invalid Authorization Header", http.StatusUnauthorized)
			return
		}

		//If token is empty Error 
		tokenString := strings.TrimSpace(parts[1])
		if tokenString == ""
		{
			http.Error(w, "empty token", http.StatusUnauthorized)
			return
		}

		
		claims := jwt.MapClaims{}

	 	token, err := jwt.ParseWithClaims(tokenString, claims, jwt.WithKeyfunc(func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok 
			{
                return nil, fmt.Errorf("unexpected signing method", t.Header["alg"])
            }
			//We should make the secret in env later
            return []byte(secret), nil
        }), 
		jwt.WithValidMethods([]string{jwt.SigningMethodS256.Alg()}),
		)


		if err != nil || !token.Valid
		{
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		//Check to make sure it is good
		var uid any
		if v, ok := claims["user_id"]; ok 
		{
			uid = v
		}
	 	else if v, ok := claims["sub"]; ok 
		{
			uid = v
		}
		else 
		{
			http.Error(w, "missing user id claim", http.StatusUnauthorized)
			return
		}


		userID, _ := uid.(string)
		if userID == "" 
		{
			http.Error(w, "invalid user id claim", http.StatusUnauthorized)
			return
		}	
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})		
}
