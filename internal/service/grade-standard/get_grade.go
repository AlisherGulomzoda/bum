package grades

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// GradesByIDs get grades by ids.
func (s Service) GradesByIDs(ctx context.Context, ids []uuid.UUID) (domain.Grades, error) {
	grades, err := s.gradesRepo.GradesBysIDTx(ctx, ids)
	if err != nil {
		return domain.Grades{}, fmt.Errorf("failed to get a grades by ids from database: %w", err)
	}

	return grades, nil
}

// GradeByID get grade by id.
func (s Service) GradeByID(ctx context.Context, id uuid.UUID) (domain.Grade, error) {
	grade, err := s.gradesRepo.GradeByIDTx(ctx, id)
	if err != nil {
		return domain.Grade{}, fmt.Errorf("failed to get a grade by grade id from database: %w", err)
	}

	return grade, nil
}
