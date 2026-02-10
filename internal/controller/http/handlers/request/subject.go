package request

import "bum-service/pkg/liblog"

// CreateSubject is a request for CreateSubject.
type CreateSubject struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description" binding:"required"`
}

// LogFields returns a list of fields for logging.
func (c CreateSubject) LogFields() liblog.Fields {
	return liblog.Fields{
		"name":        c.Name,
		"description": c.Description,
	}
}

// SubjectList is a request for SubjectList.
type SubjectList struct {
	ListFilter
}
