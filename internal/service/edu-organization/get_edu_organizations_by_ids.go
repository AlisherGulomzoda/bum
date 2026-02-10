package eduorganization

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// EduOrganizationsShortInfoByIDs get educational organizations short info by ids.
func (s Service) EduOrganizationsShortInfoByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) (domain.EduOrganizationShortInfos, error) {
	organizations, err := s.eduOrganizationRepo.EduOrganizationsShortInfoByIDsTx(ctx, ids)
	if err != nil {
		err = fmt.Errorf("failed to get educational organizations by ids from db: %w", err)
		return domain.EduOrganizationShortInfos{}, err
	}

	return organizations, nil
}
