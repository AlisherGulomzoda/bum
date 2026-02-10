package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// UpdateSchoolArgs is arguments for updating school.
type UpdateSchoolArgs struct {
	ID              uuid.UUID
	Name            string
	Location        string
	Phone           *string
	Email           *string
	GradeStandardID *uuid.UUID
}

// UpdateSchool updates school.
func (s Service) UpdateSchool(ctx context.Context, args UpdateSchoolArgs) (domain.School, error) {
	schoolDomain, err := s.schoolRepo.SchoolByIDTx(ctx, args.ID)
	if err != nil {
		return domain.School{}, fmt.Errorf("failed to get school by id: %w", err)
	}

	schoolDomain.Update(args.Name, args.Location, args.Phone, args.Email, args.GradeStandardID, s.now)
	// TODO: тут валидировать не нада как и при создании, что телефон или почта не занята гуфта?
	err = s.schoolRepo.UpdateSchoolTx(ctx, schoolDomain)
	if err != nil {
		return domain.School{}, fmt.Errorf("failed to update school to database: %w", err)
	}

	schoolDomain, err = s.SchoolByID(ctx, schoolDomain.ID)
	if err != nil {
		return domain.School{}, fmt.Errorf("failed to get school by id: %w", err)
	}

	return schoolDomain, nil
}
