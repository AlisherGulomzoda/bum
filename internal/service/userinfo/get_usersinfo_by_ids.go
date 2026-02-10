package userinfo

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// UsersByIDs get user by id.
func (s Service) UsersByIDs(ctx context.Context, ids []uuid.UUID) (domain.Users, error) {
	users, err := s.userInfoRepo.UsersByIDsTx(ctx, ids)
	if err != nil {
		return domain.Users{}, fmt.Errorf("failed to get a users by id from database: %w", err)
	}

	return users, nil
}
