package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("super-secret-key")

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateToken(tokenStr string) (int, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)

	if err != nil || !token.Valid {
		return 0, err
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	return int(userID), nil
}
