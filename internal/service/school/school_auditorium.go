package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// AuditoriumList get school subject list.
func (s Service) AuditoriumList(
	ctx context.Context,
	filters domain.AuditoriumListFilters,
) (domain.Auditoriums, int, error) {
	auditoriumList, err := s.schoolRepo.AuditoriumListTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get school auditorium list from database: %w", err)
	}

	count, err := s.schoolRepo.AuditoriumListCountTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get auditorium list count from database: %w", err)
	}

	return auditoriumList, count, nil
}

// AuditoriumByIDAndSchoolID get auditorium by id and schoolID.
func (s Service) AuditoriumByIDAndSchoolID(ctx context.Context, id, schoolID uuid.UUID) (domain.Auditorium, error) {
	auditorium, err := s.schoolRepo.AuditoriumByIDAndSchoolIDTx(ctx, id, schoolID)
	if err != nil {
		return domain.Auditorium{}, fmt.Errorf(
			"failed to get school auditorium by id and school id from database: %w", err,
		)
	}

	return auditorium, nil
}

// CreateAuditoriumArgs is a request for creating a new auditorium.
type CreateAuditoriumArgs struct {
	Name            string
	SchoolSubjectID *uuid.UUID
	Description     *string
}

// CreateAuditorium creates a new auditorium.
func (s Service) CreateAuditorium(
	ctx context.Context,
	schoolID uuid.UUID,
	req CreateAuditoriumArgs,
) (domain.Auditorium, error) {
	auditorium := domain.NewAuditorium(
		schoolID,
		req.Name,
		req.SchoolSubjectID,
		req.Description,
		s.now,
	)

	if err := s.schoolRepo.CreateAuditoriumTx(ctx, auditorium); err != nil {
		return domain.Auditorium{}, fmt.Errorf("failed to create school auditorium in database: %w", err)
	}

	return auditorium, nil
}
