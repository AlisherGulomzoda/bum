package student

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// StudentList student list.
func (s Service) StudentList(
	ctx context.Context,
	filters domain.StudentListFilter,
) (domain.Students, int, error) {
	list, err := s.studentRepo.StudentListTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get student list from database: %w", err)
	}

	var (
		userIDs   = list.UserIDs()
		schoolIDs = list.SchoolIDs()
		groupIDs  = list.GroupIDs()
	)

	users, err := s.userInfoService.UsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get user info from userinfo service: %w", err)
	}

	groups, err := s.groupService.GroupsByIDs(ctx, groupIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get groups info by ids from group service: %w", err)
	}

	schoolsShortInfo, err := s.schoolService.SchoolShortByIDs(ctx, schoolIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get user school info from service: %w", err)
	}

	list.SetUsers(users)
	list.SetShortSchools(schoolsShortInfo)
	list.SetGroups(groups)

	total, err := s.studentRepo.StudentCountTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get student list count from database: %w", err)
	}

	return list, total, nil
}
