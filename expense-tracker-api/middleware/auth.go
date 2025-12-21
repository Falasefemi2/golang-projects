package middleware

import (
	"context"
	"expense-tracker/auth"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(header, "  ")
		if len(parts) != 2 {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		userID, err := auth.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
