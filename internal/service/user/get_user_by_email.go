package user

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// UserByEmail get user by email.
func (s Service) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := s.userRepo.UserByEmailTx(ctx, email)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get the user by email from database: %w", err)
	}

	return user, nil
}
