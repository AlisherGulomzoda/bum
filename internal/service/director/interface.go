package director

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// IDirectorRepo represents a repository for director use cases.
type IDirectorRepo interface {
	AddDirectorTx(ctx context.Context, director domain.Director) error
	DirectorByIDTx(ctx context.Context, id uuid.UUID) (domain.Director, error)
	DirectorListTx(ctx context.Context, filters domain.DirectorListFilter) (domain.Directors, error)
	DirectorCountTx(ctx context.Context, filters domain.DirectorListFilter) (int, error)
}

// IUserInfoService represents a user info service for directer use cases.
type IUserInfoService interface {
	UserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UsersByIDs(ctx context.Context, ids []uuid.UUID) (domain.Users, error)
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

// ISchoolService represents a school service.
type ISchoolService interface {
	SchoolShortByIDs(ctx context.Context, ids []uuid.UUID) (domain.SchoolShortInfos, error)
}
