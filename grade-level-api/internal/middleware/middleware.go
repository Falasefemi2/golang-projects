package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/falasefemi2/gradesystem/internal/db"
	"github.com/falasefemi2/gradesystem/utils"
)

func RoleAuth(next http.HandlerFunc, requiredRole db.Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			utils.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}
		tokenString := headerParts[1]

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		user, err := db.GetUserByEmail(claims.Email)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "user not found")
			return
		}

		if db.Role(user.Role) != requiredRole {
			utils.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}