package subject

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// SubjectByID gets a subject by ID.
func (s Service) SubjectByID(
	ctx context.Context,
	id uuid.UUID,
) (subject domain.Subject, err error) {
	subject, err = s.subjectRepo.GetSubjectByIDTx(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to get subject by ID from database: %w", err)
		return domain.Subject{}, err
	}

	return subject, nil
}
