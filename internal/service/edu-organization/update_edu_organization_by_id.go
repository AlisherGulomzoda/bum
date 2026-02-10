package eduorganization

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// UpdateEduOrganizationArgs is arguments for updating Organization by id.
type UpdateEduOrganizationArgs struct {
	ID   uuid.UUID
	Name string
	Logo *string
}

// UpdateEduOrganizationByID update educational organization by id.
func (s Service) UpdateEduOrganizationByID(
	ctx context.Context,
	args UpdateEduOrganizationArgs,
) (domain.EduOrganization, error) {
	educationOrganizationEntity, err := s.eduOrganizationRepo.EduOrganizationByIDTx(ctx, args.ID)
	if err != nil {
		return domain.EduOrganization{}, fmt.Errorf("failed to get organization by id: %w", err)
	}

	educationOrganizationEntity.Update(args.Name, args.Logo, s.now)

	err = s.eduOrganizationRepo.UpdateEduOrganizationTx(ctx, educationOrganizationEntity)
	if err != nil {
		return domain.EduOrganization{}, fmt.Errorf(
			"failed to create a new educational organization to database: %w", err,
		)
	}

	return educationOrganizationEntity, nil
}
