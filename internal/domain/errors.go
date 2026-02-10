package domain

import (
	"fmt"
	"net/http"
	"strings"

	"bum-service/pkg/liberror"
)

// COMMON DEFAULT ERRORS.
var (
	// ErrInternalServerError represents an error when there is a server error.
	ErrInternalServerError = &liberror.Error{
		Err:      "Internal server error",
		Code:     "INTERNAL_SERVER_ERROR",
		HTTPCode: http.StatusInternalServerError,
	}

	// ErrBadRequest represents an error when there is a bad request.
	ErrBadRequest = &liberror.Error{
		Err:      "Bad request",
		Code:     "BAD_REQUEST",
		HTTPCode: http.StatusBadRequest,
	}

	// ErrNotFound represents an error when data was not found.
	ErrNotFound = &liberror.Error{
		Err:      "Not found",
		Code:     "NOT_FOUND",
		HTTPCode: http.StatusNotFound,
	}

	// ErrEmptyAuthHeader represents an error when auth header is empty.
	ErrEmptyAuthHeader = &liberror.Error{
		Err:      "auth header is empty",
		Code:     "AUTH_HEADER_IS_EMPTY",
		HTTPCode: http.StatusUnauthorized,
	}

	// ErrTokenIsExpired represents an error when user token is expired.
	ErrTokenIsExpired = &liberror.Error{
		Err:      "token is expired",
		Code:     "TOKEN_IS_EXPIRED",
		HTTPCode: http.StatusUnauthorized,
	}

	// ErrInvalidToken represents an error when user token is invalid.
	ErrInvalidToken = &liberror.Error{
		Err:      "invalid token",
		Code:     "INVALID_TOKEN",
		HTTPCode: http.StatusUnauthorized,
	}

	// ErrUnauthorized represents an error when user is Unauthorized.
	ErrUnauthorized = &liberror.Error{
		Err:      "Unauthorized",
		Code:     "UNAUTHORIZED",
		HTTPCode: http.StatusUnauthorized,
	}

	// ErrInvalidUser represents an error when user email or password is not correct.
	ErrInvalidUser = &liberror.Error{
		Err:      "invalid user or password",
		Code:     "INVALID_USER_OR_PASSWORD",
		HTTPCode: http.StatusUnauthorized,
	}
)

// USERS.
var (
	// ErrUserNotFound represents an error when the user not found.
	ErrUserNotFound = NewNotFoundErr("user")

	// ErrUserEmailAlreadyExists represents an error when user email is already exists.
	ErrUserEmailAlreadyExists = NewConflictErr("user email")

	// ErrUserPhoneAlreadyExists represents an error when user phone is already exists.
	ErrUserPhoneAlreadyExists = NewConflictErr("user phone")

	// ErrUserRoleBadRequest  represents an error when user role is not valid.
	ErrUserRoleBadRequest = NewBadRequest("user role")

	// ErrUserRoleInSchoolAndOrganizationAlreadyExists represents an error when user_id, school_id, organization_id
	// and role is already exists.
	ErrUserRoleInSchoolAndOrganizationAlreadyExists = NewConflictErr("user role in school and organization")

	// ErrUserGenderBadRequest  represents an error when user gender is not valid.
	ErrUserGenderBadRequest = NewBadRequest("user gender")
)

// EDUCATIONAL ORGANIZATIONS.
var (
	// ErrEduOrganizationAlreadyExists represents an error when educational organization name is already exists.
	ErrEduOrganizationAlreadyExists = NewConflictErr("educational organization name")

	// ErrEduOrganizationNotFound represents an error when educational organization is not found.
	ErrEduOrganizationNotFound = NewNotFoundErr("educational organization")
)

// SCHOOLS.
var (
	// ErrSchoolNotFound represents an error when school is not found.
	ErrSchoolNotFound = NewNotFoundErr("school")

	// ErrSchoolSubjectNotFound represents an error when school subject is not found.
	ErrSchoolSubjectNotFound = NewNotFoundErr("school subject")

	// ErrSchoolAlreadyExists represents an error when school name is already exists.
	ErrSchoolAlreadyExists = NewConflictErr("school name")

	// ErrSchoolEmailAlreadyExists represents an error when school email is already exists.
	ErrSchoolEmailAlreadyExists = NewConflictErr("school email")

	// ErrSchoolPhoneAlreadyExists represents an error when school phone is already exists.
	ErrSchoolPhoneAlreadyExists = NewConflictErr("school phone")

	// ErrHeadSchoolNotFound represents an error when school is not found.
	ErrHeadSchoolNotFound = NewNotFoundErr("head school")

	// ErrHeadquarterSchoolNotFound represents an error when school is not found.
	ErrHeadquarterSchoolNotFound = NewNotFoundErr("headquarter school")
)

// OWNERS.
var ()

// DIRECTORS.
var (
	// ErrDirectorNotFound represents an error when the director id is not found.
	ErrDirectorNotFound = NewNotFoundErr("director")
)

// HEADMASTERS.
var (
	// ErrHeadmasterNotFound represents an error when the headmaster id is not found.
	ErrHeadmasterNotFound = NewNotFoundErr("headmaster")

	// ErrHeadmasterPhoneAlreadyExists represents an error when user phone and user_id is already exists.
	ErrHeadmasterPhoneAlreadyExists = NewConflictErr("headmaster phone")

	// ErrHeadmasterEmailAlreadyExists represents an error when headmaster email and user_id is already exists.
	ErrHeadmasterEmailAlreadyExists = NewConflictErr("headmaster email")
)

// TEACHERS.
var (
	// ErrTeacherNotFound represents an error when the teacher is not found.
	ErrTeacherNotFound = NewNotFoundErr("teacher")
)

// STUDENTS.
var (
	// ErrStudentNotFound represents an error when student not found.
	ErrStudentNotFound = NewNotFoundErr("student")
)

// GUARDIANS.
var (
	// ErrGuardianAlreadyExists represents an error when guardian is already exists.
	ErrGuardianAlreadyExists = NewConflictErr("guardian")

	// ErrGuardianNotFound represents an error when guardian not found.
	ErrGuardianNotFound = NewNotFoundErr("guardian")

	// ErrStudentGuardianAlreadyExists represents an error when student guardian already exists.
	ErrStudentGuardianAlreadyExists = NewConflictErr("student guardian")

	// ErrStudentGuardianRelationBadRequest represents an error when student guardian relation is not valid.
	ErrStudentGuardianRelationBadRequest = NewBadRequest("student guardian relation")
)

// AUDITORIUMS.
var (
	// ErrAuditoriumAlreadyExists represents an error when auditorium already exists.
	ErrAuditoriumAlreadyExists = NewConflictErr("auditorium")

	// ErrAuditoriumNotFound represents an error when auditorium is not found.
	ErrAuditoriumNotFound = NewNotFoundErr("auditorium")
)

// GRADES
var (
	// ErrGradeNotFound represents an error when the grade is not found.
	ErrGradeNotFound = NewNotFoundErr("grade")

	// ErrGradeStandardNameAlreadyExists represents an error when grade standard name is already exists.
	ErrGradeStandardNameAlreadyExists = NewConflictErr("grade standard")

	// ErrGradeNameAlreadyExists represents an error when grade name is already exists.
	ErrGradeNameAlreadyExists = NewConflictErr("grade name")

	// ErrGradeStandardNotFound represents an error when the grade standard is not found.
	ErrGradeStandardNotFound = NewNotFoundErr("grade standard")
)

// GROUPS.
var (
	// ErrGroupAlreadyExists represents an error when group already exists.
	ErrGroupAlreadyExists = NewConflictErr("group")

	// ErrGroupNotFound represents an error when group is not found.
	ErrGroupNotFound = NewNotFoundErr("group")
)

// GROUP SUBJECTS.
var (
	// ErrGroupSubjectAlreadyExists represents an error when group subject already exists.
	ErrGroupSubjectAlreadyExists = NewConflictErr("group subject")

	// ErrGroupSubjectNotFound represents an error when group subject is not found.
	ErrGroupSubjectNotFound = NewNotFoundErr("group subject")
)

// STUDY PLAN.
var (
	// ErrStudyPlanAlreadyExists represents an error when study plan is already exists.
	ErrStudyPlanAlreadyExists = NewConflictErr("study plan")

	// ErrInvalidStudyPlanStatus represents an error when study plan status is not valid.
	ErrInvalidStudyPlanStatus = NewBadRequest("study plan status")
)

// SUBJECTS.
var (
	// ErrSubjectNameAlreadyExists represents an error when subject name is already exists.
	ErrSubjectNameAlreadyExists = NewConflictErr("subject name")

	// ErrSubjectNotFound represents an error when subject is not found.
	ErrSubjectNotFound = NewNotFoundErr("subject")
)

// LESSONS.
var (
	// ErrLessonNotFound represents an error when lesson is not found.
	ErrLessonNotFound = NewNotFoundErr("lesson")
)

// MARKS.
var (
	// ErrMarkAlreadyExists represents an error when mark name is already exists.
	ErrMarkAlreadyExists = NewConflictErr("mark")
)

// NewNotFoundErr creates a new NotFound error with the given entity.
func NewNotFoundErr(entity string) *liberror.Error {
	return &liberror.Error{
		Err:      fmt.Sprintf("%s not found", entity),
		Code:     fmt.Sprintf("NOT_FOUND: %s", strings.ToUpper(entity)),
		HTTPCode: http.StatusNotFound,
	}
}

// NewConflictErr creates a new Conflict error with the given entity.
func NewConflictErr(entity string) *liberror.Error {
	return &liberror.Error{
		Err:      fmt.Sprintf("%s already exists", entity),
		Code:     fmt.Sprintf("CONFLICT: %s", strings.ToUpper(entity)),
		HTTPCode: http.StatusConflict,
	}
}

// NewBadRequest creates a new bad request error with the given message.
func NewBadRequest(message string) *liberror.Error {
	return &liberror.Error{
		Err:      message,
		Code:     "BAD_REQUEST",
		HTTPCode: http.StatusBadRequest,
	}
}
