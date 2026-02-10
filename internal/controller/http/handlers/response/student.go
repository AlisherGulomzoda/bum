package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// Student is student response.
type Student struct {
	User   User      `json:"user"`
	ID     uuid.UUID `json:"id"`
	RoleID uuid.UUID `json:"role_id"`
	UserID uuid.UUID `json:"user_id"`

	GroupID uuid.UUID `json:"group_id"`
	Group   Group     `json:"group"`

	SchoolID        uuid.UUID        `json:"school_id"`
	SchoolShortInfo *SchoolShortInfo `json:"school_short_info,omitempty"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewStudent creates a new student response.
func NewStudent(student domain.Student) Student {
	return Student{
		User:   NewUser(student.User),
		ID:     student.ID,
		RoleID: student.RoleID,
		UserID: student.UserID,

		GroupID: student.GroupID,
		Group:   NewGroup(student.Group),

		SchoolID:        student.SchoolID,
		SchoolShortInfo: NewSchoolShortInfo(student.SchoolShortInfo),

		CreatedAt: utils.RFC3339Time(student.CreatedAt),
		UpdatedAt: utils.RFC3339Time(student.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(student.DeletedAt),
	}
}

// StudentShortInfo is student short info response.
type StudentShortInfo struct {
	User   User      `json:"user"`
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`

	GroupID  uuid.UUID `json:"group_id"`
	SchoolID uuid.UUID `json:"school_id"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewStudentShortInfo creates a new student short info response.
func NewStudentShortInfo(student *domain.Student) *StudentShortInfo {
	if student == nil {
		return nil
	}

	return &StudentShortInfo{
		ID:       student.ID,
		UserID:   student.UserID,
		GroupID:  student.GroupID,
		SchoolID: student.SchoolID,

		User: NewUser(student.User),

		CreatedAt: utils.RFC3339Time(student.CreatedAt),
		UpdatedAt: utils.RFC3339Time(student.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(student.DeletedAt),
	}
}

// StudentList response model for listing students.
type StudentList struct {
	Students   []Student  `json:"students"`
	Pagination Pagination `json:"pagination"`
}

// NewStudentList creates a new student list for response.
func NewStudentList(
	list domain.Students,
	pagination Pagination,
) StudentList {
	students := make([]Student, len(list))

	for i := range list {
		students[i] = NewStudent(list[i])
	}

	return StudentList{
		Students:   students,
		Pagination: pagination,
	}
}
