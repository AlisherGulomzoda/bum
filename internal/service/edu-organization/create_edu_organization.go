package eduorganization

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// CreateEduOrganizationArgs is arguments for creating a new Organization.
type CreateEduOrganizationArgs struct {
	Name        string
	Logo        *string
	Description *string
}

// CreateEduOrganization create new educational organization.
func (s Service) CreateEduOrganization(
	ctx context.Context,
	args CreateEduOrganizationArgs,
) (domain.EduOrganization, error) {
	educationOrganizationEntity := domain.NewEduOrganization(args.Name, args.Logo, args.Description, s.now)

	err := s.eduOrganizationRepo.CreateEduOrganizationTx(ctx, educationOrganizationEntity)
	if err != nil {
		return domain.EduOrganization{}, fmt.Errorf(
			"failed to create a new educational organization to database: %w", err,
		)
	}

	return educationOrganizationEntity, nil
}
