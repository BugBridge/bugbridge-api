package handlers

import 
(
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type User struct
{
	ID string
	Email string
	PasswordHash string
}


//Mongo Stuff


type AuthService struct 
{
	Secret   []byte
	Issuer 	 string
	Audience string
	TTL 	 time.Duration
}


func NewAuthServiceFromEnv() *AuthService 
{
	return &AuthService{
		Secret:   []byte(os.Getenv("SECRET")),
		Issuer:   "bugbridge-api",
		Audience: "bugbridge-frontend",
		TTL:      2 * time.Hour,
	}
}

func (a *AuthService) Sign(userID string) (string, error)
{
	if len(a.secret) == 0
	{
		return "", errors.New("Missing secret")
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": a.Issuer,
		"aud": a.Audience,
		"iat": time.Now().unix,
		"exp": time.Now().Add(a.TTL).Unix,	
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJwtToken := jwtToken.SignedString(a.secret)
	if err != nil
	{
		return "", err
	}

	return signedJwtToken, nil
}

