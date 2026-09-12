package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"cookapp/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Token invalide", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Token invalide", http.StatusUnauthorized)
			return
		}
		userID := int(claims["id"].(float64))
		ctx := context.WithValue(r.Context(), model.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
