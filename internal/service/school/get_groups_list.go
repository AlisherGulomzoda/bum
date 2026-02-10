package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// GroupList get group list.
func (s Service) GroupList(
	ctx context.Context,
	schoolID uuid.UUID,
	filters domain.GroupFilters,
) (domain.Groups, int, error) {
	groupList, err := s.groupRepo.GroupListTx(ctx, schoolID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get group list from database: %w", err)
	}

	count, err := s.groupRepo.GroupListCountTx(ctx, schoolID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get school list count from database: %w", err)
	}

	gradeIDs := groupList.GradeIDs()

	grades, err := s.gradeService.GradesByIDs(ctx, gradeIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get grades by ids: %w", err)
	}

	groupList.SetGrades(grades)

	return groupList, count, nil
}
