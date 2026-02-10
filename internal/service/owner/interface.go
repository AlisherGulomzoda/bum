package owner

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

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

// IOwnerRepository represents a owner repository for owner use cases.
type IOwnerRepository interface {
	AddOwnerTx(ctx context.Context, owner domain.Owner) error
	OwnerByIDTx(ctx context.Context, id uuid.UUID) (domain.Owner, error)
	OwnerByUserIDAndSchoolIDTx(ctx context.Context, schoolID, userID uuid.UUID) (domain.Owner, error)
	OwnerListTx(ctx context.Context, filters domain.OwnerListFilter) (domain.Owners, error)
	OwnerCountTx(ctx context.Context, filters domain.OwnerListFilter) (int, error)
}

// IUserInfoService represents a user info service for directer use cases.
type IUserInfoService interface {
	UserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UsersByIDs(ctx context.Context, ids []uuid.UUID) (domain.Users, error)
}
