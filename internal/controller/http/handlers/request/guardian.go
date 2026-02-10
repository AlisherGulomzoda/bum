package request

import (
	"github.com/google/uuid"

	"bum-service/pkg/liblog"
)

// AssignStudentGuardian is assign student guardian request.
type AssignStudentGuardian struct {
	UserID   uuid.UUID `json:"user_id" binding:"required,uuid"`
	SchoolID uuid.UUID `json:"school_id" binding:"required,uuid"`
	Relation string    `json:"relation" binding:"required,oneof=mother father guardian relative"`
}

// LogFields returns a list of fields for logging.
func (c AssignStudentGuardian) LogFields() liblog.Fields {
	return liblog.Fields{
		"user_id":   c.UserID,
		"school_id": c.SchoolID,
		"relation":  c.Relation,
	}
}

// StudentGuardianList request model for listing of student guardians.
type StudentGuardianList struct {
	ListFilter

	GroupIDsFilter
	SchoolIDsFilter
	OrganizationIDsFilter

	CreatedDate DateFilter
}
