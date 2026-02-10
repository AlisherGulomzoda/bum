package grades

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// CreateGradeStandardArgs is arguments for creating a new grade standard.
type CreateGradeStandardArgs struct {
	OrganizationID *uuid.UUID
	Name           string
	EducationYears int8
	Description    *string
	Grades         []CreateGradeArgs
}

// CreateGradeArgs is arguments for creating a new grade.
type CreateGradeArgs struct {
	Name          string
	EducationYear *int8
}

// CreateGradeStandard creates a new grade-standard.
func (s Service) CreateGradeStandard(
	ctx context.Context,
	arg CreateGradeStandardArgs,
) (domain.GradeStandard, error) {
	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.GradeStandard{}, fmt.Errorf("failed to begin transaction : %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf(
				"failed to end transaction on create grade standard: %w: %w",
				domain.ErrInternalServerError,
				errEnd,
			)
		}
	}(tx)

	newGradeStandard := domain.NewGradeStandard(
		arg.OrganizationID,
		arg.Name,
		arg.EducationYears,
		arg.Description,
		s.now,
	)

	err = s.gradesRepo.CreateGradeStandardTx(txCtx, newGradeStandard)
	if err != nil {
		return domain.GradeStandard{}, fmt.Errorf("failed create grade standard to database: %w", err)
	}

	newGrades := createGradesDomain(arg, newGradeStandard.ID, s.now)

	err = s.gradesRepo.CreateGradesTx(txCtx, newGrades)
	if err != nil {
		return domain.GradeStandard{}, fmt.Errorf("failed create grades to database: %w", err)
	}

	savedGradeStandard, err := s.GradeStandardByID(txCtx, newGradeStandard.ID)
	if err != nil {
		return domain.GradeStandard{}, fmt.Errorf("failed to get a grade standard by id: %w", err)
	}

	return savedGradeStandard, nil
}

func createGradesDomain(
	args CreateGradeStandardArgs,
	gradeStandardID uuid.UUID,
	nowFunc func() time.Time,
) domain.Grades {
	grades := make(domain.Grades, len(args.Grades))

	for i, grade := range args.Grades {
		grades[i] = domain.NewGrade(gradeStandardID, grade.Name, grade.EducationYear, nowFunc)
	}

	return grades
}
