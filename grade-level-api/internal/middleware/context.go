// middleware/context.go
package middleware

import (
	"context"
	"errors"

	"github.com/falasefemi2/gradesystem/internal/models"
)

func GetUserFromContext(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value("user").(*models.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
