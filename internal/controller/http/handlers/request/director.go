package request

import (
	"github.com/google/uuid"

	"bum-service/pkg/liblog"
)

// AddDirector is a request to create a new director.
type AddDirector struct {
	UserID   uuid.UUID `json:"user_id" binding:"required,uuid"`
	SchoolID uuid.UUID `json:"school_id" binding:"required,uuid"`
	Phone    *string   `json:"phone" binding:"omitempty,e164"`
	Email    *string   `json:"email" binding:"omitempty,email"`
}

// LogFields returns a list of fields for logging.
func (c AddDirector) LogFields() liblog.Fields {
	return liblog.Fields{
		"user_id":   c.UserID,
		"school_id": c.SchoolID,
		"phone":     c.Phone,
		"email":     c.Email,
	}
}

// DirectorList request model for listing of directors.
type DirectorList struct {
	ListFilter

	SchoolIDsFilter

	CreatedDate DateFilter
}
