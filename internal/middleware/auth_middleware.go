package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/teodora-lazarevic/Poll-App/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// RequireAuth extracts, validates JWT token and injects userID into request context.
func RequireAuth(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorJSON(writer, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.ErrorJSON(writer, http.StatusUnauthorized, "Invalid token format")
			return
		}

		// Validate token with security check for signing algorithm
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			utils.ErrorJSON(writer, http.StatusUnauthorized, "Unauthorized or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.ErrorJSON(writer, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			utils.ErrorJSON(writer, http.StatusUnauthorized, "Invalid sub claim")
			return
		}

		// Inject User ID into context
		ctx := context.WithValue(request.Context(), UserIDKey, int(sub))
		next(writer, request.WithContext(ctx), ps)
	}
}

// GetUserIDFromContext retrieves the authenticated user ID safely.
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}
