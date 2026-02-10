package eduorganization

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// IEduOrganizationRepo represents a repository for educational organization use cases.
type IEduOrganizationRepo interface {
	CreateEduOrganizationTx(ctx context.Context, organization domain.EduOrganization) error
	UpdateEduOrganizationTx(ctx context.Context, o domain.EduOrganization) error
	EduOrganizationByIDTx(ctx context.Context, id uuid.UUID) (domain.EduOrganization, error)
	EduOrganizationsByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.EduOrganizations, error)
	EduOrganizationsShortInfoByIDsTx(
		ctx context.Context, ids []uuid.UUID,
	) (domain.EduOrganizationShortInfos, error)
	EduOrganizationShortByIDTx(
		ctx context.Context, id uuid.UUID,
	) (domain.EduOrganizationShortInfo, error)
	EduOrganizationListTx(ctx context.Context, filters domain.EduOrganizationFilters) (domain.EduOrganizations, error)
	EduOrganizationCountTx(ctx context.Context) (int, error)
}
