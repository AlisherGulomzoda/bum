package repository

import (
	"database/sql"
	"errors"
	"strings"

	"bum-service/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
)

// whereAnd complements the selection for the query.
func where(parameters []string) string {
	if len(parameters) == 0 {
		return ""
	}

	return " WHERE " + strings.Join(parameters, " AND ")
}

//nolint:gochecknoglobals //it's map of errors
var mapOfErrors = map[string]error{
	// Director errors
	DirectorsUserIDFKey:   domain.ErrUserNotFound,
	DirectorsSchoolIDFKey: domain.ErrSchoolNotFound,

	// Educational organization errors
	EduOrganizationsNameKey: domain.ErrEduOrganizationAlreadyExists,

	// Owner errors
	OwnersUserIDFKey:         domain.ErrUserNotFound,
	OwnersOrganizationIDFKey: domain.ErrEduOrganizationNotFound,

	// Headmaster errors
	HeadmastersUserIDFKey:   domain.ErrUserNotFound,
	HeadmastersSchoolIDFKey: domain.ErrSchoolNotFound,

	// Schools errors
	SchoolsNameKey:            domain.ErrSchoolAlreadyExists,
	SchoolsOrganizationIDFKey: domain.ErrEduOrganizationNotFound,
	SchoolGradeStandardIDFKey: domain.ErrGradeStandardNotFound,

	// School subjects
	SchoolSubjectsSchoolIDFKey:  domain.ErrSchoolNotFound,
	SchoolSubjectsSubjectIDFKey: domain.ErrSubjectNotFound,

	// Subjects errors
	SubjectsNameUniqueKey: domain.ErrSubjectNameAlreadyExists,

	// Users errors
	UsersEmailUniqueKey: domain.ErrUserEmailAlreadyExists,
	UsersPhoneUniqueKey: domain.ErrUserPhoneAlreadyExists,

	// User roles errors
	UserRolesUserIDRoleSchoolIDOrganizationIDUniqueKey: domain.ErrUserRoleInSchoolAndOrganizationAlreadyExists,
	UserRolesUserIDFKey:         domain.ErrUserNotFound,
	UserRolesSchoolIDFKey:       domain.ErrSchoolNotFound,
	UserRolesOrganizationIDFKey: domain.ErrEduOrganizationNotFound,

	// Teacher errors
	TeachersUserIDFKey:   domain.ErrUserNotFound,
	TeachersSchoolIDFKey: domain.ErrSchoolNotFound,

	// Grade errors
	GradesNameUniqueKey:       domain.ErrGradeNameAlreadyExists,
	GradesGradeStandardIDFKey: domain.ErrGradeStandardNotFound,

	// Grade Standard errors
	GradeStandardsNameUniqueKey:      domain.ErrGradeStandardNameAlreadyExists,
	GradeStandardsOrganizationIDFKey: domain.ErrEduOrganizationNotFound,

	// Group errors
	GroupsGradeIDFkey:          domain.ErrGradeNotFound,
	GroupNameKey:               domain.ErrGroupAlreadyExists,
	GroupsSchoolIDFkey:         domain.ErrSchoolNotFound,
	GroupsClassTeacherIDFKey:   domain.ErrTeacherNotFound,
	GroupsClassPresidentIDFKey: domain.ErrStudentNotFound,

	// Group subjects errors
	GroupSubjectsUniqueKey:           domain.ErrGroupSubjectAlreadyExists,
	GroupSubjectsSchoolSubjectIDFkey: domain.ErrSchoolSubjectNotFound,
	GroupSubjectsGroupIDFkey:         domain.ErrGroupNotFound,
	GroupSubjectsTeacherIDFKey:       domain.ErrTeacherNotFound,

	// Study Plan errors
	StudyPlanGroupSubjectPlanOrderUniqueKey: domain.ErrStudyPlanAlreadyExists,
	StudyPlanGroupSubjectIDFKey:             domain.ErrGroupSubjectNotFound,

	// Auditorium errors
	AuditoriumsSchoolSubjectIDFkey: domain.ErrSchoolSubjectNotFound,
	AuditoriumsSchoolIDFkey:        domain.ErrSchoolNotFound,
	AuditoriumsNameUniqueKey:       domain.ErrAuditoriumAlreadyExists,

	// Student errors
	StudentsGroupIDFKey: domain.ErrGroupNotFound,
	StudentsUserIDFKey:  domain.ErrUserNotFound,

	// Lesson errors
	LessonsSchoolIDFKey:       domain.ErrSchoolNotFound,
	LessonsGroupSubjectIDFKey: domain.ErrGroupSubjectNotFound,
	LessonsTeacherIDFKey:      domain.ErrTeacherNotFound,
	LessonsAuditoriumIDFKey:   domain.ErrAuditoriumNotFound,

	// Marks
	MarksLessonIDFKey:  domain.ErrLessonNotFound,
	MarksStudentIDFKey: domain.ErrStudentNotFound,
	MarksLessonKey:     domain.ErrMarkAlreadyExists,

	// Student Guardians
	StudentGuardiansKey:           domain.ErrStudentGuardianAlreadyExists,
	StudentGuardiansStudentIDFKey: domain.ErrStudentNotFound,
	StudentGuardiansUserIDFKey:    domain.ErrUserNotFound,
	StudentGuardiansSchoolIDFKey:  domain.ErrSchoolNotFound,
}

func handleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if handledError, ok := mapOfErrors[pgErr.ConstraintName]; ok {
			return handledError
		}

		return pgErr
	}

	return err
}
