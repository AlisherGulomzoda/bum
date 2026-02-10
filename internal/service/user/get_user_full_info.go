package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// UserFullInfoByID get user full info by id.
func (s Service) UserFullInfoByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.userRepo.UserByIDTx(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get the user by id from database: %w", err)
	}

	roles, err := s.UserRoles(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get the user roles by id: %w", err)
	}

	user.SetRoles(roles)

	return user, nil
}
