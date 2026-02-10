package subject

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// SubjectList gets a subject list.
func (s Service) SubjectList(
	ctx context.Context,
	filters domain.SubjectListFilter,
) (subjects []domain.Subject, count int, err error) {
	subjects, err = s.subjectRepo.GetSubjectListTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get subject list from database: %w", err)
	}

	count, err = s.subjectRepo.SubjectCountTx(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get subject count from database: %w", err)
	}

	return subjects, count, nil
}
