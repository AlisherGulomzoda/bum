package headmaster

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// HeadmasterList headmaster list.
func (s Service) HeadmasterList(
	ctx context.Context,
	filters domain.HeadmasterListFilter,
) ([]domain.Headmaster, int, error) {
	list, err := s.headmasterRepo.HeadmasterListTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get headmaster list database: %w", err)
	}

	var (
		userIDs   = list.UserIDs()
		schoolIDs = list.SchoolIDs()
	)

	users, err := s.userInfoService.UsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get user info from userinfo service: %w", err)
	}

	schoolsShortInfo, err := s.schoolService.SchoolShortByIDs(ctx, schoolIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get user school info from service: %w", err)
	}

	list.SetUsers(users)
	list.SetShortSchools(schoolsShortInfo)

	total, err := s.headmasterRepo.HeadmasterCountTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get headmaster list count database: %w", err)
	}

	return list, total, nil
}
