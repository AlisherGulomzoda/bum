package grades

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// GradeStandardList grade standard list.
func (s Service) GradeStandardList(
	ctx context.Context,
	filter domain.GradeStandardListFilter,
) (domain.GradeStandards, int, error) {
	list, err := s.gradesRepo.GradeStandardListTx(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get grade standard list database: %w", err)
	}

	total, err := s.gradesRepo.GradeStandardListCountTx(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get grade standard list count database: %w", err)
	}

	return list, total, nil
}
