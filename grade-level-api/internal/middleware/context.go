package middleware

import (
	"context"

	"github.com/falasefemi2/gradesystem/internal/models"
)

type contextKey string

const userContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
