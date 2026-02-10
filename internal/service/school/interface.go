package school

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// ISchoolRepo represents a repository for school use cases.
//
//nolint:interfacebloat // it's ok
type ISchoolRepo interface {
	CreateSchoolTx(ctx context.Context, organization domain.School) error
	SchoolByIDTx(ctx context.Context, schoolID uuid.UUID) (domain.School, error)
	UpdateSchoolTx(ctx context.Context, o domain.School) error
	SchoolListTx(ctx context.Context, filters domain.SchoolFilters) (domain.Schools, error)
	SchoolListCountTx(ctx context.Context, filters domain.SchoolFilters) (int, error)
	SchoolShortByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.SchoolShortInfos, error)
	SchoolShortByIDTx(ctx context.Context, id uuid.UUID) (domain.SchoolShortInfo, error)

	CreateSchoolSubjectTx(ctx context.Context, o domain.SchoolSubject) error
	SchoolSubjectByIDAndSchoolIDTx(
		ctx context.Context,
		id uuid.UUID,
		schoolID uuid.UUID,
	) (domain.SchoolSubject, error)
	SchoolSubjectsByIDs(ctx context.Context, ids []uuid.UUID) (domain.SchoolSubjects, error)
	SchoolSubjectListTx(ctx context.Context, filters domain.SchoolSubjectFilters) (domain.SchoolSubjects, error)
	SchoolSubjectsListCountTx(ctx context.Context, filters domain.SchoolSubjectFilters) (int, error)

	CreateAuditoriumTx(ctx context.Context, o domain.Auditorium) error
	AuditoriumByIDAndSchoolIDTx(ctx context.Context, id, schoolID uuid.UUID) (domain.Auditorium, error)
	AuditoriumListTx(ctx context.Context, filters domain.AuditoriumListFilters) (domain.Auditoriums, error)
	AuditoriumListCountTx(ctx context.Context, filters domain.AuditoriumListFilters) (int, error)

	AssignStudyPlansTx(ctx context.Context, groupSubjectID uuid.UUID, studyPlans domain.StudyPlans) error
	StudyPlanListTx(ctx context.Context, groupSubjectID uuid.UUID) (domain.StudyPlans, error)
	StudyPlanChangeStatusTx(ctx context.Context, groupSubjectID, studyPlanID uuid.UUID, status string) error
}

// IGroupSubjectsRepo represents a repository for group subjects.
type IGroupSubjectsRepo interface {
	AddGroupSubjectsTx(ctx context.Context, groupSubject domain.GroupSubject) error
	GroupSubjectByIDTx(ctx context.Context, id uuid.UUID) (domain.GroupSubject, error)
	GroupSubjectListTx(ctx context.Context, groupID uuid.UUID) (domain.GroupSubjects, error)
}

// IGroupRepo is a repository for groups.
type IGroupRepo interface {
	CreateGroupTx(ctx context.Context, g domain.Group) error
	GroupByIDTx(ctx context.Context, id uuid.UUID) (domain.Group, error)
	UpdateGroupTx(ctx context.Context, group domain.Group) error
	GroupsByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.Groups, error)
	GroupListTx(
		ctx context.Context, schoolID uuid.UUID, filters domain.GroupFilters,
	) (domain.Groups, error)
	GroupListCountTx(
		ctx context.Context, schoolID uuid.UUID, filters domain.GroupFilters,
	) (int, error)
}

// IEduOrganizationService represents an edu organization service.
type IEduOrganizationService interface {
	EduOrganizationShortByID(ctx context.Context, id uuid.UUID) (domain.EduOrganizationShortInfo, error)
	EduOrganizationsShortInfoByIDs(ctx context.Context, ids []uuid.UUID) (domain.EduOrganizationShortInfos, error)
}

// IGradeService represents grade service.
type IGradeService interface {
	GradesByIDs(ctx context.Context, ids []uuid.UUID) (domain.Grades, error)
	GradeByID(ctx context.Context, id uuid.UUID) (domain.Grade, error)
}

// ITeacherService represents teacher service.
type ITeacherService interface {
	TeachersByIDs(ctx context.Context, ids []uuid.UUID) (domain.Teachers, error)
	TeacherByID(ctx context.Context, id uuid.UUID) (domain.Teacher, error)
}

// IStudentService represents student service.
type IStudentService interface {
	StudentShortInfoByID(ctx context.Context, studentID uuid.UUID) (domain.Student, error)
}
