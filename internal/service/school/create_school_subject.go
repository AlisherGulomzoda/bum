package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// CreateSchoolSubjectArgs is a list of arguments to create a new school subject.
type CreateSchoolSubjectArgs struct {
	SchoolID    uuid.UUID
	SubjectID   uuid.UUID
	Name        string
	Description *string
}

// CreateSchoolSubject creates a new school subject.
func (s Service) CreateSchoolSubject(ctx context.Context, args CreateSchoolSubjectArgs) (domain.SchoolSubject, error) {
	schoolSubjectDomain := domain.NewSchoolSubject(
		args.SchoolID,
		args.SubjectID,
		args.Name,
		args.Description,
		s.now,
	)

	err := s.schoolRepo.CreateSchoolSubjectTx(ctx, schoolSubjectDomain)
	if err != nil {
		return domain.SchoolSubject{}, fmt.Errorf("failed to create a new school subject to database: %w", err)
	}

	schoolSubjectDomain, err = s.SchoolSubjectByIDAndSchoolID(ctx, schoolSubjectDomain.ID, args.SchoolID)
	if err != nil {
		return domain.SchoolSubject{}, fmt.Errorf("failed to get school by id: %w", err)
	}

	return schoolSubjectDomain, nil
}
