package handlers

import (
	"context"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/internal/service/director"
	eduorganization "bum-service/internal/service/edu-organization"
	grades "bum-service/internal/service/grade-standard"
	"bum-service/internal/service/headmaster"
	"bum-service/internal/service/lesson"
	"bum-service/internal/service/owner"
	"bum-service/internal/service/school"
	"bum-service/internal/service/student"
	"bum-service/internal/service/subject"
	"bum-service/internal/service/teacher"
	"bum-service/internal/service/user"
)

// ISystemService is a System use case interface.
type ISystemService interface {
	HealthCheck(ctx context.Context) error
}

// IEduOrganizationService is an educational organization use case interface.
//
//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=interface.go -destination=mocks/edu_organization_service.go
type IEduOrganizationService interface {
	CreateEduOrganization(
		ctx context.Context,
		args eduorganization.CreateEduOrganizationArgs,
	) (domain.EduOrganization, error)
	UpdateEduOrganizationByID(
		ctx context.Context,
		args eduorganization.UpdateEduOrganizationArgs,
	) (domain.EduOrganization, error)
	EduOrganizationByID(ctx context.Context, id uuid.UUID) (domain.EduOrganization, error)
	EduOrganizationList(
		ctx context.Context, filters domain.EduOrganizationFilters,
	) (domain.EduOrganizations, int, error)
}

// ISchoolService is a school use case interface.
//
//nolint:interfacebloat // it's ok
type ISchoolService interface {
	AddSchool(ctx context.Context, args school.CreateSchoolArgs) (domain.School, error)
	SchoolByID(ctx context.Context, schoolID uuid.UUID) (domain.School, error)
	UpdateSchool(ctx context.Context, args school.UpdateSchoolArgs) (domain.School, error)
	SchoolList(ctx context.Context, filters domain.SchoolFilters) (domain.Schools, int, error)

	CreateGroup(ctx context.Context, arg school.CreateGroupArgs) (domain.Group, error)
	GroupByID(ctx context.Context, groupID uuid.UUID) (domain.Group, error)
	UpdateGroup(ctx context.Context, args school.UpdateGroupArgs) (domain.Group, error)
	GroupList(ctx context.Context, schoolID uuid.UUID, filters domain.GroupFilters) (domain.Groups, int, error)

	CreateSchoolSubject(ctx context.Context, args school.CreateSchoolSubjectArgs) (domain.SchoolSubject, error)
	SchoolSubjectByIDAndSchoolID(
		ctx context.Context,
		id uuid.UUID,
		schoolID uuid.UUID,
	) (domain.SchoolSubject, error)
	SchoolSubjectList(
		ctx context.Context,
		filters domain.SchoolSubjectFilters,
	) (domain.SchoolSubjects, int, error)

	AddGroupSubject(
		ctx context.Context, schoolID, groupID uuid.UUID, args school.AddGroupSubjectArgs,
	) (domain.GroupSubject, error)
	GroupSubjectList(ctx context.Context, groupID uuid.UUID) (domain.GroupSubjects, error)

	CreateAuditorium(
		ctx context.Context,
		schoolID uuid.UUID,
		args school.CreateAuditoriumArgs,
	) (domain.Auditorium, error)
	AuditoriumByIDAndSchoolID(ctx context.Context, id, schoolID uuid.UUID) (domain.Auditorium, error)
	AuditoriumList(ctx context.Context, filters domain.AuditoriumListFilters) (domain.Auditoriums, int, error)

	AssignStudyPlans(ctx context.Context, schoolID uuid.UUID, args []school.AddStudyPlanArgs) (domain.StudyPlans, error)
	StudyPlanList(ctx context.Context, groupSubjectID uuid.UUID) (domain.StudyPlans, error)
	StudyPlanChangeStatus(ctx context.Context, groupSubjectID, studyPlanID uuid.UUID, status string) error
}

// IDirectorService is a director use case interface.
type IDirectorService interface {
	AddDirector(ctx context.Context, args director.AddDirectorArgs) (domain.Director, error)
	DirectorByID(ctx context.Context, id uuid.UUID) (domain.Director, error)
	DirectorList(ctx context.Context, filters domain.DirectorListFilter) ([]domain.Director, int, error)
}

// IHeadmasterService is a headmaster use case interface.
type IHeadmasterService interface {
	AddHeadmaster(ctx context.Context, args headmaster.AddHeadmasterArgs) (domain.Headmaster, error)
	HeadmasterByID(ctx context.Context, id uuid.UUID) (domain.Headmaster, error)
	HeadmasterList(ctx context.Context, filters domain.HeadmasterListFilter) ([]domain.Headmaster, int, error)
}

// ISubjectService is a subject use case interface.
type ISubjectService interface {
	CreateSubject(ctx context.Context, args subject.CreateSubjectArgs) (domain.Subject, error)
	SubjectByID(ctx context.Context, id uuid.UUID) (domain.Subject, error)
	SubjectList(ctx context.Context, filters domain.SubjectListFilter) ([]domain.Subject, int, error)
}

// IUserService is a user use case interface.
type IUserService interface {
	UserList(ctx context.Context, filters domain.UserListFilter) (domain.Users, int, error)
	UserRoles(ctx context.Context, id uuid.UUID) (domain.UserRoles, error)
	UserFullInfoByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	AddUser(ctx context.Context, args user.AddUserArgs) (newUser domain.User, err error)
}

// ITeacherService is a teacher use case interface.
type ITeacherService interface {
	AddTeacher(ctx context.Context, args teacher.AddTeacherArgs) (domain.Teacher, error)
	TeacherByID(ctx context.Context, id uuid.UUID) (domain.Teacher, error)
	TeacherList(ctx context.Context, filters domain.TeacherListFilter) (domain.Teachers, int, error)
}

// IGradesService is a grades use case interface.
type IGradesService interface {
	CreateGradeStandard(ctx context.Context, arg grades.CreateGradeStandardArgs) (domain.GradeStandard, error)
	GradeStandardByID(ctx context.Context, id uuid.UUID) (domain.GradeStandard, error)
	GradeStandardList(ctx context.Context, filter domain.GradeStandardListFilter) (domain.GradeStandards, int, error)
}

// IAuthService is an auth service interface.
type IAuthService interface {
	GetUserIDByEmailAndPassword(ctx context.Context, email, password string) (uuid.UUID, error)
}

// IStudentService is student service interface.
type IStudentService interface {
	AddStudent(ctx context.Context, args student.AddStudentArgs) (newStudent domain.Student, err error)
	StudentByID(ctx context.Context, studentID uuid.UUID) (domain.Student, error)
	StudentList(ctx context.Context, filters domain.StudentListFilter) (domain.Students, int, error)

	StudentGuardians(ctx context.Context, studentID uuid.UUID) (domain.StudentGuardians, error)
	AssignStudentGuardian(ctx context.Context, args student.AssignStudentGuardianArgs) (domain.StudentGuardian, error)
	StudentGuardianByUserID(ctx context.Context, userID uuid.UUID) (domain.Guardian, error)
	StudentGuardianList(
		ctx context.Context,
		filters domain.StudentGuardianListFilter,
	) (domain.StudentGuardians, int, error)
}

// ILessonService is lesson service interface.
type ILessonService interface {
	AssignLessons(ctx context.Context, args lesson.AddWeekLessonsArgs) (domain.Lessons, error)
	LessonsList(ctx context.Context, filters domain.LessonsListFilter) (domain.Lessons, error)

	AddMark(ctx context.Context, args lesson.AddMarkArgs) (domain.Mark, error)
	MarkByID(ctx context.Context, markID uuid.UUID) (domain.Mark, error)
}

// IOwnerService is owner service interface.
type IOwnerService interface {
	AddOwner(ctx context.Context, arg owner.AddOwnerArgs) (newOwner domain.Owner, err error)
	OwnerByID(ctx context.Context, id uuid.UUID) (domain.Owner, error)
	OwnerByUserIDAndSchoolID(ctx context.Context, schoolID, userID uuid.UUID) (domain.Owner, error)
	OwnerList(
		ctx context.Context,
		filter domain.OwnerListFilter,
	) (domain.Owners, int, error)
}
