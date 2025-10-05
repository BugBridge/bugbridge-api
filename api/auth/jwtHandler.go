package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	Secret   []byte
	Issuer   string
	Audience string
	TTL      time.Duration
}

// Maybe be changed to get values from config/config.go
func NewAuthServiceFromEnv() *AuthService {
	return &AuthService{
		Secret:   []byte(os.Getenv("SECRET")),
		Issuer:   "bugbridge-api",
		Audience: "bugbridge-frontend",
		TTL:      2 * time.Hour,
	}
}

func (a *AuthService) Sign(userID string) (string, error) {
	if len(a.Secret) == 0 {
		return "", errors.New("missing secret")
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"iss": a.Issuer,
		"aud": a.Audience,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(a.TTL).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJwtToken, err := jwtToken.SignedString(a.Secret)

	if err != nil {
		return "", err
	}

	return signedJwtToken, nil
}

// Parses the tokens
func (a *AuthService) Parse(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	//Parse claims e.g userID, (users Id) returns 2 things
	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return a.Secret, nil
		},

		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(a.Issuer),
		jwt.WithAudience(a.Audience),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(30*time.Second),
	)

	if err != nil {
		return nil, err
	}

	return claims, nil
}

// Returns the subject (user Id) from jwt token
func Subject(claims jwt.MapClaims) (string, bool) {
	//is sub there
	if v, ok := claims["sub"]; ok {
		//if sub is a string good else bad
		if s, ok := v.(string); ok && s != "" {
			return s, true
		}

		//if v is float convert to string int no decimal
		if f, ok := v.(float64); ok {
			return strconv.FormatInt(int64(f), 10), true
		}
	}

	return "", false
}
