package eduorganization

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// EduOrganizationByID get educational organization by id.
func (s Service) EduOrganizationByID(ctx context.Context, id uuid.UUID) (domain.EduOrganization, error) {
	createdOrganization, err := s.eduOrganizationRepo.EduOrganizationByIDTx(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to create a new educational organization to database: %w", err)
		return domain.EduOrganization{}, err
	}

	return createdOrganization, nil
}

// EduOrganizationShortByID get educational organization short info by id.
func (s Service) EduOrganizationShortByID(ctx context.Context, id uuid.UUID) (domain.EduOrganizationShortInfo, error) {
	createdOrganization, err := s.eduOrganizationRepo.EduOrganizationShortByIDTx(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to create a new educational organization to database: %w", err)
		return domain.EduOrganizationShortInfo{}, err
	}

	return createdOrganization, nil
}
