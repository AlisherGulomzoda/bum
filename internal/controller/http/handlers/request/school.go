package request

import (
	"github.com/google/uuid"

	"bum-service/pkg/liblog"
)

// AddSchool is a request to add a new school.
type AddSchool struct {
	Name           string  `json:"name" binding:"required"`
	OrganizationID string  `json:"organization_id" binding:"required,uuid"`
	Location       string  `json:"location" binding:"required"`
	Phone          *string `json:"phone" binding:"omitempty,e164"`
	Email          *string `json:"email" binding:"omitempty,email"`
}

// OrganizationUUID converts OrganizationID to uuid.
func (c AddSchool) OrganizationUUID() uuid.UUID {
	return uuid.MustParse(c.OrganizationID)
}

// LogFields returns a list of fields for logging.
func (c AddSchool) LogFields() liblog.Fields {
	return liblog.Fields{
		"name":            c.Name,
		"organization_id": c.OrganizationID,
		"location":        c.Location,
		"phone":           c.Phone,
		"email":           c.Email,
	}
}

// SchoolList is a request to get school list.
type SchoolList struct {
	ListFilter

	Phones []string `form:"phones[]" binding:"omitempty,dive,e164"`
	Emails []string `form:"emails[]" binding:"omitempty,dive,email"`

	OrganizationIDsFilter
}

// GroupList is a request to get group list.
type GroupList struct {
	ListFilter
}

// CreateSchoolSubject is a request to create a new subject for school.
type CreateSchoolSubject struct {
	SubjectID   uuid.UUID `json:"subject_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description" binding:"required"`
}

// LogFields returns a list of fields for logging.
func (c CreateSchoolSubject) LogFields() liblog.Fields {
	return liblog.Fields{
		"name":       c.Name,
		"subject_id": c.SubjectID,
	}
}

// SchoolSubjectList is a request to get school subject list.
type SchoolSubjectList struct {
	ListFilter
}

// AssignStudyPlan is a request to assign study plan to group subject.
type AssignStudyPlan struct {
	ID          *uuid.UUID `json:"id,omitempty" binding:"omitnil,uuid"`
	Title       string     `json:"title" binding:"required"`
	Description *string    `json:"description,omitempty" binding:"omitnil"`
}

// UpdateSchool is request for updating school.
type UpdateSchool struct {
	Name     string  `json:"name" binding:"required"`
	Location string  `json:"location" binding:"required"`
	Phone    *string `json:"phone" binding:"omitempty,e164"`
	Email    *string `json:"email" binding:"omitempty,email"`

	GradeStandardID *uuid.UUID `json:"grade_standard_id,omitempty" binding:"omitempty,uuid"`
}

// SchoolIDsFilter school ids filter for embedding.
type SchoolIDsFilter struct {
	SchoolIDs []string `form:"school_ids[]" binding:"omitempty,dive,uuid"`
}

// SchoolUUIDs converts SchoolIDs to list of uuid.
func (s SchoolIDsFilter) SchoolUUIDs() []uuid.UUID {
	if len(s.SchoolIDs) == 0 {
		return []uuid.UUID{}
	}

	uuids := make([]uuid.UUID, len(s.SchoolIDs))
	for idx := range s.SchoolIDs {
		uuids[idx] = uuid.MustParse(s.SchoolIDs[idx])
	}

	return uuids
}
