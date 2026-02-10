package request

import (
	"github.com/google/uuid"

	"bum-service/pkg/liblog"
)

// CreateAuditorium is a request for CreateAuditorium.
type CreateAuditorium struct {
	Name             string     `json:"name" binding:"required"`
	SchoolSubjectsID *uuid.UUID `json:"school_subject_id" binding:"omitnil"`
	Description      *string    `json:"description" binding:"omitnil"`
}

// LogFields returns a list of fields for logging.
func (c CreateAuditorium) LogFields() liblog.Fields {
	return liblog.Fields{
		"name":               c.Name,
		"school_subjects_id": c.SchoolSubjectsID,
		"description":        c.Description,
	}
}

// AuditoriumList is a request to get AuditoriumList.
type AuditoriumList struct {
	ListFilter
}
