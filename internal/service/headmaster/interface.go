package headmaster

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// IHeadmasterRepo represents a repository for headmaster use cases.
type IHeadmasterRepo interface {
	AddHeadmasterTx(ctx context.Context, headmaster domain.Headmaster) error
	HeadmasterByIDTx(ctx context.Context, id uuid.UUID) (domain.Headmaster, error)
	HeadmasterListTx(ctx context.Context, filters domain.HeadmasterListFilter) (domain.Headmasters, error)
	HeadmasterCountTx(ctx context.Context, filters domain.HeadmasterListFilter) (int, error)
}

// IUserInfoService represents a user info service for headmaster use cases.
type IUserInfoService interface {
	UserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UsersByIDs(ctx context.Context, ids []uuid.UUID) (domain.Users, error)
}

// ISchoolService represents a school service.
type ISchoolService interface {
	SchoolShortByIDs(ctx context.Context, ids []uuid.UUID) (domain.SchoolShortInfos, error)
}

// IUserService represents a user service for adding roles.
type IUserService interface {
	AddRoleToUser(
		ctx context.Context,
		userID uuid.UUID,
		role domain.Role,
		schoolID *uuid.UUID,
		organizationID *uuid.UUID,
	) (newUserRole domain.UserRole, err error)
}
