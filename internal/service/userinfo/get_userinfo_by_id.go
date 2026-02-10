package userinfo

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// UserByID get user by id.
func (s Service) UserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.userInfoRepo.UserByIDTx(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get a user by id from database: %w", err)
	}

	return user, nil
}
