package school

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// SchoolList get school list.
func (s Service) SchoolList(
	ctx context.Context,
	filters domain.SchoolFilters,
) (domain.Schools, int, error) {
	schoolList, err := s.schoolRepo.SchoolListTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get school list from database: %w", err)
	}

	organizationIDs := schoolList.GetOrganizationIDs()

	organizations, err := s.eduOrganizationService.EduOrganizationsShortInfoByIDs(ctx, organizationIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get organizations by ids: %w", err)
	}

	schoolList.SetOrganizationShortInfos(organizations)

	count, err := s.schoolRepo.SchoolListCountTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get school list count from database: %w", err)
	}

	return schoolList, count, nil
}
