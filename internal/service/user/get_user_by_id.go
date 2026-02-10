package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// UserByID get user by id.
func (s Service) UserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	user, err := s.userRepo.UserByIDTx(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get the user by id from database: %w", err)
	}

	return user, nil
}
