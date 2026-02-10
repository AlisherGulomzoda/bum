package owner

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// OwnerList owner list.
func (s Service) OwnerList(
	ctx context.Context,
	filter domain.OwnerListFilter,
) (domain.Owners, int, error) {
	list, err := s.ownerRepo.OwnerListTx(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get owners list database: %w", err)
	}

	userIDs := list.UserIDs()

	users, err := s.userInfoService.UsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get user info from userinfo service: %w", err)
	}

	list.SetUsers(users)

	total, err := s.ownerRepo.OwnerCountTx(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get owners list count database: %w", err)
	}

	return list, total, nil
}
