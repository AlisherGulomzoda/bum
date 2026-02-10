package request

import (
	"github.com/google/uuid"

	"bum-service/pkg/liblog"
)

// AddTeacher is request for create teacher.
type AddTeacher struct {
	UserID   uuid.UUID `json:"user_id" binding:"required,uuid"`
	SchoolID uuid.UUID `json:"school_id" binding:"required,uuid"`
	Phone    *string   `json:"phone" binding:"omitempty,e164"`
	Email    *string   `json:"email" binding:"omitempty,email"`
}

// LogFields returns a list of fields for logging.
func (a AddTeacher) LogFields() liblog.Fields {
	return liblog.Fields{
		"user_id":   a.UserID,
		"phone":     a.Phone,
		"email":     a.Email,
		"school_id": a.SchoolID,
	}
}

// TeacherList request model for listing of teachers.
type TeacherList struct {
	ListFilter

	SchoolIDsFilter
	GroupIDsFilter
	OrganizationIDsFilter

	CreatedDate DateFilter
}
