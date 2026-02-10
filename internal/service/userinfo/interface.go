package userinfo

import (
	"context"

	"bum-service/internal/domain"

	"github.com/google/uuid"
)

// IUserInfoRepo represents a user repository for user use cases.
type IUserInfoRepo interface {
	UserByIDTx(ctx context.Context, id uuid.UUID) (domain.User, error)
	UsersByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.Users, error)
}
