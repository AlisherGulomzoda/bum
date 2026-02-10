package grades

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// IGradesRepo represents a repository for grade-standard use cases.
type IGradesRepo interface {
	CreateGradeStandardTx(ctx context.Context, gradeStandard domain.GradeStandard) error
	CreateGradesTx(ctx context.Context, grades domain.Grades) error

	GradeStandardByIDTx(ctx context.Context, id uuid.UUID) (domain.GradeStandard, error)
	GradesByGradeStandardIDTx(ctx context.Context, gradeStandardID uuid.UUID) (domain.Grades, error)
	GradesBysIDTx(ctx context.Context, ids []uuid.UUID) (domain.Grades, error)
	GradeByIDTx(ctx context.Context, id uuid.UUID) (domain.Grade, error)

	GradeStandardListTx(ctx context.Context, filters domain.GradeStandardListFilter) (domain.GradeStandards, error)
	GradeStandardListCountTx(ctx context.Context, filters domain.GradeStandardListFilter) (int, error)
}
