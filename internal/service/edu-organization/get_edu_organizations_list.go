package eduorganization

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// EduOrganizationList get a list of educational organizations.
func (s Service) EduOrganizationList(
	ctx context.Context,
	filters domain.EduOrganizationFilters,
) (domain.EduOrganizations, int, error) {
	list, err := s.eduOrganizationRepo.EduOrganizationListTx(ctx, filters)
	if err != nil {
		err = fmt.Errorf("failed to get educational organization list to database: %w", err)
		return nil, 0, err
	}

	count, err := s.eduOrganizationRepo.EduOrganizationCountTx(ctx)
	if err != nil {
		err = fmt.Errorf("failed to get educational organization count to database: %w", err)
		return nil, 0, err
	}

	return list, count, nil
}
