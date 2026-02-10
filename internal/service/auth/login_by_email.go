package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// GetUserIDByEmailAndPassword returns userID by email and password.
func (s Service) GetUserIDByEmailAndPassword(ctx context.Context, email, password string) (uuid.UUID, error) {
	// Get user by email in order to get its id and password
	user, err := s.userService.UserByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to check password: %w : %w", domain.ErrInvalidUser, err)
	}

	// Check the password
	if err = utils.ComparePassword(user.Password, password); err != nil {
		return uuid.Nil, fmt.Errorf("failed to check password: %w : %w", domain.ErrInvalidUser, err)
	}

	return user.ID, nil
}
