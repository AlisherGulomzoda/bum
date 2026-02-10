package grades

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// GradeStandardByID get grade standard by id.
func (s Service) GradeStandardByID(
	ctx context.Context,
	id uuid.UUID,
) (domain.GradeStandard, error) {
	gradeStandard, err := s.gradesRepo.GradeStandardByIDTx(ctx, id)
	if err != nil {
		return domain.GradeStandard{}, fmt.Errorf("failed to get a grade standard by id from database: %w", err)
	}

	grades, err := s.gradesRepo.GradesByGradeStandardIDTx(ctx, id)
	if err != nil {
		return domain.GradeStandard{}, fmt.Errorf("failed to get grades by grade standard id from database: %w", err)
	}

	gradeStandard.SetGrades(grades)

	return gradeStandard, nil
}
