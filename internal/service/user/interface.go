package user

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// IUserRepo represents a user repository for user use cases.
type IUserRepo interface {
	AddUserTx(ctx context.Context, user domain.User) error
	AddUserRoleTx(ctx context.Context, role domain.UserRole, withinOnConflict bool) error
	UserByIDTx(ctx context.Context, id uuid.UUID) (domain.User, error)
	UserByEmailTx(ctx context.Context, email string) (domain.User, error)
	GetUserListTx(ctx context.Context, filters domain.UserListFilter) (domain.Users, error)
	UserCountTx(ctx context.Context, filters domain.UserListFilter) (int, error)

	UserRolesByIDTx(ctx context.Context, userID uuid.UUID) (domain.UserRoles, error)
	UserRolesByIDsTx(ctx context.Context, userIDs []uuid.UUID) (domain.UserRoles, error)
}

// ISchoolService represents a school service for users.
type ISchoolService interface {
	SchoolShortByIDs(ctx context.Context, ids []uuid.UUID) (domain.SchoolShortInfos, error)
}

// IStudentRepo represents student repo.
// TODO убрать после рефакторинга.
type IStudentRepo interface {
	AddStudentTx(ctx context.Context, o domain.Student) error
	StudentByIDTx(ctx context.Context, id uuid.UUID) (domain.Student, error)
	StudentsByUserIDTx(ctx context.Context, userID uuid.UUID) (domain.Students, error)
	StudentsByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.Students, error)
	StudentListTx(ctx context.Context, filters domain.StudentListFilter) (domain.Students, error)
	StudentCountTx(ctx context.Context, filters domain.StudentListFilter) (int, error)

	AddStudentGuardianTx(ctx context.Context, studentGuardian domain.StudentGuardian) error
	StudentGuardianByIDTx(ctx context.Context, id uuid.UUID) (domain.StudentGuardian, error)
	StudentGuardiansByStudentIDTx(ctx context.Context, studentID uuid.UUID) (domain.StudentGuardians, error)
	StudentGuardianListTx(ctx context.Context, filters domain.StudentGuardianListFilter) (domain.StudentGuardians, error)
	StudentGuardianListCountTx(ctx context.Context, filters domain.StudentGuardianListFilter) (int, error)
}

// IGroupService represents group service.
type IGroupService interface {
	GroupByID(ctx context.Context, groupID uuid.UUID) (domain.Group, error)
	GroupsByIDs(ctx context.Context, ids []uuid.UUID) (domain.Groups, error)
}

// IEduOrganizationService represents an edu organization service.
type IEduOrganizationService interface {
	EduOrganizationByID(ctx context.Context, id uuid.UUID) (domain.EduOrganization, error)
	EduOrganizationsShortInfoByIDs(ctx context.Context, ids []uuid.UUID) (domain.EduOrganizationShortInfos, error)
}
