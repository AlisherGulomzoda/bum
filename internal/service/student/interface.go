package student

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// IStudentRepo represents student repo.
type IStudentRepo interface {
	AddStudentTx(ctx context.Context, o domain.Student) error
	StudentByIDTx(ctx context.Context, id uuid.UUID) (domain.Student, error)
	StudentsByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.Students, error)
	StudentListTx(ctx context.Context, filters domain.StudentListFilter) (domain.Students, error)
	StudentCountTx(ctx context.Context, filters domain.StudentListFilter) (int, error)

	AddStudentGuardianTx(ctx context.Context, studentGuardian domain.StudentGuardian) error
	StudentGuardianByIDTx(ctx context.Context, id uuid.UUID) (domain.StudentGuardian, error)
	StudentGuardiansByStudentIDTx(ctx context.Context, studentID uuid.UUID) (domain.StudentGuardians, error)
	StudentGuardianByUserIDTx(ctx context.Context, id uuid.UUID) (domain.StudentGuardians, error)
	StudentGuardianListTx(ctx context.Context, filters domain.StudentGuardianListFilter) (domain.StudentGuardians, error)
	StudentGuardianListCountTx(ctx context.Context, filters domain.StudentGuardianListFilter) (int, error)
}

// IUserInfoService represents a user info service for headmaster use cases.
type IUserInfoService interface {
	UserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UsersByIDs(ctx context.Context, ids []uuid.UUID) (domain.Users, error)
}

// ISchoolService represents a school service.
type ISchoolService interface {
	SchoolShortByIDs(ctx context.Context, ids []uuid.UUID) (domain.SchoolShortInfos, error)
	SchoolShortByID(ctx context.Context, id uuid.UUID) (domain.SchoolShortInfo, error)
}

// IGroupService represents group service.
type IGroupService interface {
	GroupByID(ctx context.Context, groupID uuid.UUID) (domain.Group, error)
	GroupsByIDs(ctx context.Context, ids []uuid.UUID) (domain.Groups, error)
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
