package director

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// DirectorList director list.
func (s Service) DirectorList(
	ctx context.Context,
	filter domain.DirectorListFilter,
) ([]domain.Director, int, error) {
	list, err := s.directorRepo.DirectorListTx(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get director list database: %w", err)
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

	total, err := s.directorRepo.DirectorCountTx(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get director list count database: %w", err)
	}

	return list, total, nil
}
