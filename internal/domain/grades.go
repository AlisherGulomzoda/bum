package domain

import (
	"time"

	"github.com/google/uuid"
)

// GradeStandard is structure of grade standard.
type GradeStandard struct {
	ID             uuid.UUID
	OrganizationID *uuid.UUID
	Organization   *EduOrganization
	Name           string
	EducationYears int8
	Description    *string
	Grades         Grades

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewGradeStandard creates a new GradeStandard domain.
func NewGradeStandard(
	organizationID *uuid.UUID,
	name string,
	educationYears int8,
	description *string,
	nowFunc func() time.Time,
) GradeStandard {
	now := nowFunc()

	return GradeStandard{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		Name:           name,
		EducationYears: educationYears,
		Description:    description,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// SetGrades set grade standard grades.
func (gs *GradeStandard) SetGrades(grades Grades) {
	gs.Grades = grades
}

// GradeStandards are collection of GradeStandard.
type GradeStandards []GradeStandard

// Grade is structure of grade.
type Grade struct {
	ID              uuid.UUID
	GradeStandardID uuid.UUID
	Name            string
	EducationYear   *int8
	Groups          *Groups

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewGrade creates a new Grade domain.
func NewGrade(
	gradeStandardID uuid.UUID,
	name string,
	educationYear *int8,
	nowFunc func() time.Time,
) Grade {
	now := nowFunc()

	return Grade{
		ID:              uuid.New(),
		GradeStandardID: gradeStandardID,
		Name:            name,
		EducationYear:   educationYear,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Grades are collection of Grade.
type Grades []Grade

// SetGroups sets groups for grades.
func (g Grades) SetGroups(groups Groups) {
	mapOfGroupsByGradeID := make(map[uuid.UUID]Groups)
	for _, group := range groups {
		mapOfGroupsByGradeID[group.GradeID] = append(mapOfGroupsByGradeID[group.GradeID], group)
	}

	for index := range g {
		gradeGroups := mapOfGroupsByGradeID[g[index].ID]
		g[index].Groups = &gradeGroups
	}
}

// GradeStandardListFilter filter for the list of grade standard.
type GradeStandardListFilter struct {
	ListFilter
}

// NewGradeStandardListFilter creates a new GradeStandardListFilter domain.
func NewGradeStandardListFilter(list ListFilter) GradeStandardListFilter {
	return GradeStandardListFilter{
		ListFilter: list,
	}
}
