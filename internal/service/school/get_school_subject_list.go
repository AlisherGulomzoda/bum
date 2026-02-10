package school

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// SchoolSubjectList get school subject list.
func (s Service) SchoolSubjectList(
	ctx context.Context,
	filters domain.SchoolSubjectFilters,
) (domain.SchoolSubjects, int, error) {
	// TODO: добавить проверку на существование school id.
	schoolSubjectList, err := s.schoolRepo.SchoolSubjectListTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get school subject list from database: %w", err)
	}

	count, err := s.schoolRepo.SchoolSubjectsListCountTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get school list count from database: %w", err)
	}

	return schoolSubjectList, count, nil
}
