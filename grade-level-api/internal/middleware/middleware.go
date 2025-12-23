package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/falasefemi2/gradesystem/internal/db"
	"github.com/falasefemi2/gradesystem/utils"
)

func RoleAuth(next http.HandlerFunc, allowedRoles ...db.Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		claims, err := ValidateJWT(parts[1])
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		user, err := db.GetUserByEmail(claims.Email)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "user not found")
			return
		}

		authorized := false
		for _, role := range allowedRoles {
			if db.Role(user.Role) == role {
				authorized = true
				break
			}
		}

		if !authorized {
			utils.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
