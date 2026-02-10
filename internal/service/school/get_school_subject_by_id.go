package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// SchoolSubjectByIDAndSchoolID get school by id and schoolID.
func (s Service) SchoolSubjectByIDAndSchoolID(
	ctx context.Context,
	id uuid.UUID,
	schoolID uuid.UUID,
) (domain.SchoolSubject, error) {
	schoolSubject, err := s.schoolRepo.SchoolSubjectByIDAndSchoolIDTx(ctx, id, schoolID)
	if err != nil {
		return domain.SchoolSubject{}, fmt.Errorf(
			"failed to get school subject by id and school id from database: %w", err,
		)
	}

	return schoolSubject, nil
}
