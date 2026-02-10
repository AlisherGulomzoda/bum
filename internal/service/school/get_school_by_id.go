package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// SchoolByID get school by id.
func (s Service) SchoolByID(ctx context.Context, schoolID uuid.UUID) (domain.School, error) {
	school, err := s.schoolRepo.SchoolByIDTx(ctx, schoolID)
	if err != nil {
		return domain.School{}, fmt.Errorf("failed to get a school by id from database: %w", err)
	}

	organization, err := s.eduOrganizationService.EduOrganizationShortByID(ctx, school.OrganizationID)
	if err != nil {
		return domain.School{}, fmt.Errorf("failed to get organization by id: %w", err)
	}

	school.SetOrganizationShortInfo(organization)

	return school, nil
}
