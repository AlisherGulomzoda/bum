package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// GradeStandard is a structure of grade standard response.
type GradeStandard struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Name           string     `json:"name"`
	EducationYears int8       `json:"education_years"`
	Description    *string    `json:"description"`
	Grades         []Grade    `json:"grades"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// Grade is a structure of grade response.
type Grade struct {
	ID              string  `json:"id"`
	GradeStandardID string  `json:"grade_standard_id"`
	Name            string  `json:"name"`
	EducationYear   *int8   `json:"education_year,omitempty"`
	Groups          *Groups `json:"groups,omitempty"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewGradeStandard creates a new grade standard response.
func NewGradeStandard(gs domain.GradeStandard) GradeStandard {
	grades := make([]Grade, len(gs.Grades))

	for i, grade := range gs.Grades {
		grades[i] = NewGrade(grade)
	}

	return GradeStandard{
		ID:             gs.ID,
		OrganizationID: gs.OrganizationID,
		Name:           gs.Name,
		EducationYears: gs.EducationYears,
		Description:    gs.Description,
		Grades:         grades,

		CreatedAt: utils.RFC3339Time(gs.CreatedAt),
		UpdatedAt: utils.RFC3339Time(gs.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(gs.DeletedAt),
	}
}

// NewGrade creates a new grade response.
func NewGrade(g domain.Grade) Grade {
	var groups *Groups

	if g.Groups != nil {
		gList := NewGroups(*g.Groups)
		groups = &gList
	}

	return Grade{
		ID:              g.ID.String(),
		GradeStandardID: g.GradeStandardID.String(),
		Name:            g.Name,
		EducationYear:   g.EducationYear,
		Groups:          groups,

		CreatedAt: utils.RFC3339Time(g.CreatedAt),
		UpdatedAt: utils.RFC3339Time(g.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(g.DeletedAt),
	}
}

// GradeStandardList response model for listing grade standard.
type GradeStandardList struct {
	GradeStandards []GradeStandard `json:"grade_standards"`
	Pagination     Pagination      `json:"pagination"`
}

// NewGradeStandardList creates a new grade standard list for response.
func NewGradeStandardList(
	list domain.GradeStandards,
	pagination Pagination,
) GradeStandardList {
	gradeStandards := make([]GradeStandard, len(list))

	for i := range list {
		gradeStandards[i] = NewGradeStandard(list[i])
	}

	return GradeStandardList{
		GradeStandards: gradeStandards,
		Pagination:     pagination,
	}
}

// Grades is a list of grade.
type Grades []Grade

// GroupsList is list response for groups.
type GroupsList struct {
	Groups Groups `json:"groups"`
	Total  int    `json:"total"`
}

// NewGroupsList creates a new grades list response from domain grades.
func NewGroupsList(list domain.Groups, total int) GroupsList {
	return GroupsList{
		Groups: NewGroups(list),
		Total:  total,
	}
}
