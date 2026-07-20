package handlers

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Validates a JWT token
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// Decoding token and verifying signature
	return jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
}

// Gets the authenticated user ID from the request
func GetAuthenticatedUserID(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")

	// Removing prefix from the token string
	tokenString := strings.TrimPrefix(
		authHeader,
		"Bearer ",
	)

	token, err := ValidateJWT(tokenString)
	if err != nil {
		return 0, err
	}

	// .(jwt.MapClaims) is used to convert the claims into a map - something like type assertion
	claims, ok := token.Claims.(jwt.MapClaims) // Converting into map
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// .(float64) is used to convert the value into a float64 - something like type assertion
	return int(claims["sub"].(float64)), nil
}

// Generates a JWT token (signed with secret key)
func GenerateJWTToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET"))) // Signing token with secret key
}
