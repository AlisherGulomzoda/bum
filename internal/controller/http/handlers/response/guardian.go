package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// StudentGuardian is student guardian response.
type StudentGuardian struct {
	ID        uuid.UUID `json:"id"`
	StudentID uuid.UUID `json:"student_id"`
	Student   Student   `json:"student"`
	UserID    uuid.UUID `json:"user_id"`
	User      User      `json:"user"`
	SchoolID  uuid.UUID `json:"school_id"`
	Relation  string    `json:"relation"`

	CreatedAt utils.RFC3339Time `json:"created_at"`
	UpdatedAt utils.RFC3339Time `json:"updated_at"`
}

// NewStudentGuardian creates a new student guardian response.
func NewStudentGuardian(studentGuardian domain.StudentGuardian) StudentGuardian {
	return StudentGuardian{
		ID:        studentGuardian.ID,
		StudentID: studentGuardian.StudentID,
		Student:   NewStudent(studentGuardian.Student),
		UserID:    studentGuardian.UserID,
		User:      NewUser(studentGuardian.User),
		SchoolID:  studentGuardian.SchoolID,
		Relation:  string(studentGuardian.Relation),

		CreatedAt: utils.RFC3339Time(studentGuardian.CreatedAt),
		UpdatedAt: utils.RFC3339Time(studentGuardian.UpdatedAt),
	}
}

// StudentGuardians is student guardians response.
type StudentGuardians []StudentGuardian

// NewStudentGuardians creates a new student guardians response.
func NewStudentGuardians(studentGuardians domain.StudentGuardians) StudentGuardians {
	guardians := make(StudentGuardians, len(studentGuardians))
	for index, studentGuardian := range studentGuardians {
		guardians[index] = NewStudentGuardian(studentGuardian)
	}

	return guardians
}

// StudentGuardianList response model for listing student guardians.
type StudentGuardianList struct {
	StudentGuardians StudentGuardians `json:"guardians"`
	Pagination       Pagination       `json:"pagination"`
}

// NewStudentGuardianList creates a new student guardian list for response.
func NewStudentGuardianList(
	list domain.StudentGuardians,
	pagination Pagination,
) StudentGuardianList {
	return StudentGuardianList{
		StudentGuardians: NewStudentGuardians(list),
		Pagination:       pagination,
	}
}

// StudentGuardiansForStudent is student guardians for student response.
type StudentGuardiansForStudent struct {
	StudentGuardians StudentGuardians `json:"guardians"`
}

// NewStudentGuardiansForStudent creates a new student guardians for student response.
func NewStudentGuardiansForStudent(studentGuardians domain.StudentGuardians) StudentGuardiansForStudent {
	return StudentGuardiansForStudent{
		StudentGuardians: NewStudentGuardians(studentGuardians),
	}
}

type GuardianStudent struct {
	ID       uuid.UUID `json:"id"`
	SchoolID uuid.UUID `json:"school_id"`
	Relation string    `json:"relation"`

	StudentID uuid.UUID `json:"student_id"`
	Student   Student   `json:"student"`

	CreatedAt utils.RFC3339Time `json:"created_at"`
	UpdatedAt utils.RFC3339Time `json:"updated_at"`
}

func NewGuardianStudent(studentGuardian domain.StudentGuardian) GuardianStudent {
	return GuardianStudent{
		ID:       studentGuardian.ID,
		SchoolID: studentGuardian.SchoolID,

		Relation: string(studentGuardian.Relation),

		StudentID: studentGuardian.StudentID,
		Student:   NewStudent(studentGuardian.Student),

		CreatedAt: utils.RFC3339Time(studentGuardian.CreatedAt),
		UpdatedAt: utils.RFC3339Time(studentGuardian.UpdatedAt),
	}
}

func NewGuardianStudents(studentGuardians domain.StudentGuardians) GuardianStudents {
	res := make(GuardianStudents, 0, len(studentGuardians))
	for _, studentGuardian := range studentGuardians {
		res = append(res, NewGuardianStudent(studentGuardian))
	}

	return res
}

type GuardianStudents []GuardianStudent

type Guardian struct {
	User             User `json:"user"`
	GuardianStudents GuardianStudents
}

// NewGuardian creates a new guardian response.
func NewGuardian(guardian domain.Guardian) Guardian {
	return Guardian{
		User:             NewUser(guardian.User),
		GuardianStudents: NewGuardianStudents(guardian.StudentGuardians),
	}
}
