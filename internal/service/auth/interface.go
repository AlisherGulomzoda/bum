package auth

import (
	"context"

	"bum-service/internal/domain"
)

// IUserService represents a user service for adding roles.
type IUserService interface {
	UserByEmail(ctx context.Context, email string) (domain.User, error)
}
