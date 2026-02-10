package teacher

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// ITeacherRepo represents a repository for teacher use cases.
type ITeacherRepo interface {
	CreateTeacherTx(ctx context.Context, o domain.Teacher) error
	TeacherByIDTx(ctx context.Context, id uuid.UUID) (domain.Teacher, error)
	TeachersByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.Teachers, error)
	TeacherListTx(ctx context.Context, filters domain.TeacherListFilter) (domain.Teachers, error)
	TeacherCountTx(ctx context.Context, filters domain.TeacherListFilter) (int, error)
}

// IUserInfoService represents a user info service for teacher use cases.
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
