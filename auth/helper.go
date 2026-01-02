package main

import (
	"errors"
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

func ValidateJWT(token string, key []byte) (string, error) {

	parsedToken, err := jwt.Parse(
		token,
		func(t *jwt.Token) (any, error) {
			return key, nil
		},
	)
	if err != nil || !parsedToken.Valid {
		return "", err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims type")
	}

	email, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("missing or invalid user email in token")
	}

	return email, nil
}
