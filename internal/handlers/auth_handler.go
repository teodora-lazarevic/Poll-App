package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
}

func GetAuthenticatedUserID(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")

	tokenString := strings.TrimPrefix(
		authHeader,
		"Bearer ",
	)

	token, err := ValidateJWT(tokenString)
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return int(claims["sub"].(float64)), nil
}

func authenticate(writer http.ResponseWriter, request *http.Request) (int, bool) {
	userId, err := GetAuthenticatedUserID(request)
	if err != nil {
		ErrorJSON(writer, http.StatusUnauthorized, "Unauthorized")
		return 0, false
	}
	return userId, true
}

// Generates a JWT token
func GenerateJWTToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
