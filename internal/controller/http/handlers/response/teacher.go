package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// Teacher is a structure of teacher response.
type Teacher struct {
	ID              uuid.UUID        `json:"id"`
	RoleID          uuid.UUID        `json:"role_id"`
	User            User             `json:"user"`
	SchoolID        uuid.UUID        `json:"school_id"`
	SchoolShortInfo *SchoolShortInfo `json:"school_short_info,omitempty"`
	Phone           *string          `json:"phone,omitempty"`
	Email           *string          `json:"email,omitempty"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// TeacherShortInfo is teacher short information.
type TeacherShortInfo struct {
	ID    uuid.UUID `json:"id"`
	User  User      `json:"user"`
	Phone *string   `json:"phone,omitempty"`
	Email *string   `json:"email,omitempty"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewTeacherShortInfo returns a new teacher short info response.
func NewTeacherShortInfo(teacher *domain.Teacher) *TeacherShortInfo {
	if teacher == nil {
		return nil
	}

	return &TeacherShortInfo{
		ID:    teacher.ID,
		User:  NewUser(teacher.User),
		Phone: teacher.Phone,
		Email: teacher.Email,

		CreatedAt: utils.RFC3339Time(teacher.CreatedAt),
		UpdatedAt: utils.RFC3339Time(teacher.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(teacher.DeletedAt),
	}
}

// NewTeacher creates a new teacher response.
func NewTeacher(teacher domain.Teacher) Teacher {
	return Teacher{
		ID:              teacher.ID,
		RoleID:          teacher.RoleID,
		User:            NewUser(teacher.User),
		SchoolID:        teacher.SchoolID,
		SchoolShortInfo: NewSchoolShortInfo(teacher.SchoolShortInfo),
		Phone:           teacher.Phone,
		Email:           teacher.Email,

		CreatedAt: utils.RFC3339Time(teacher.CreatedAt),
		UpdatedAt: utils.RFC3339Time(teacher.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(teacher.DeletedAt),
	}
}

// TeacherList response model for listing teachers.
type TeacherList struct {
	Teachers   []Teacher  `json:"teachers"`
	Pagination Pagination `json:"pagination"`
}

// NewTeacherList creates a new teacher list for response.
func NewTeacherList(
	list []domain.Teacher,
	pagination Pagination,
) TeacherList {
	teachers := make([]Teacher, len(list))

	for i := range list {
		teachers[i] = NewTeacher(list[i])
	}

	return TeacherList{
		Teachers:   teachers,
		Pagination: pagination,
	}
}
