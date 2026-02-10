package teacher

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// TeacherList teacher list.
func (s Service) TeacherList(ctx context.Context, filters domain.TeacherListFilter) (domain.Teachers, int, error) {
	list, err := s.teacherRepo.TeacherListTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get teacher list database: %w", err)
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

	total, err := s.teacherRepo.TeacherCountTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get teacher list count database: %w", err)
	}

	return list, total, nil
}
