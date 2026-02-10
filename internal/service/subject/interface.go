package subject

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// ISubjectRepo represents a repository for subject use cases.
type ISubjectRepo interface {
	CreateSubjectTx(ctx context.Context, subject domain.Subject) error
	GetSubjectByIDTx(ctx context.Context, id uuid.UUID) (domain.Subject, error)
	GetSubjectListTx(ctx context.Context, filters domain.SubjectListFilter) (domain.Subjects, error)
	SubjectCountTx(ctx context.Context) (int, error)
}
