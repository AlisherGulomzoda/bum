package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// SchoolSubjectByIDs get school subjects by ids.
func (s Service) SchoolSubjectByIDs(ctx context.Context, ids []uuid.UUID) (domain.SchoolSubjects, error) {
	schoolSubjectList, err := s.schoolRepo.SchoolSubjectsByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get school subjects by ids from database: %w", err)
	}

	return schoolSubjectList, nil
}
