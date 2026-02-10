package user

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// UserList gets a user list.
func (s Service) UserList(
	ctx context.Context, filters domain.UserListFilter,
) (users domain.Users, count int, err error) {
	users, err = s.userRepo.GetUserListTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user list from database: %w", err)
	}

	userIDs := users.UserIDs()

	userRoles, err := s.UserRolesByIDs(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user roles from service: %w", err)
	}

	users.SetRoles(userRoles)

	count, err = s.userRepo.UserCountTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user count from database: %w", err)
	}

	return users, count, nil
}
