package main

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(email string) (string, error) {
	key := []byte(os.Getenv("SIGNING_KEY"))
	now := time.Now()
	
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "auth-service",
			"sub": email,
			"iat": now.Unix(),
			"exp": now.Add(24 * time.Hour).Unix(),
		},
	)
	
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	
	return signedToken, nil
}