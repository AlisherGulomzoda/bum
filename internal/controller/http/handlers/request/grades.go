package request

import (
	"github.com/google/uuid"

	"bum-service/pkg/liblog"
)

// CreateGradeStandard is a request to create a new grade standard.
type CreateGradeStandard struct {
	OrganizationID *uuid.UUID    `json:"organization_id"`
	Name           string        `json:"name" binding:"required"`
	EducationYears int8          `json:"education_years" binding:"required"`
	Description    *string       `json:"description" binding:"required"`
	Grades         []CreateGrade `json:"grades" binding:"required,min=1"`
}

// CreateGrade is a request to create a new grade.
type CreateGrade struct {
	Name          string `json:"name" binding:"required"`
	EducationYear *int8  `json:"education_year" binding:"omitempty"`
}

// LogFields returns a list of fields for logging.
func (c CreateGradeStandard) LogFields() liblog.Fields {
	return liblog.Fields{
		"organization_id": c.OrganizationID,
		"name":            c.Name,
		"education_years": c.EducationYears,
		"description":     c.Description,
		"grades":          c.Grades,
	}
}

// GradeStandardList request model for listing of grade standard.
type GradeStandardList struct {
	ListFilter
}
