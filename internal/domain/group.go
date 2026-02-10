package domain

import (
	"time"

	"github.com/google/uuid"
)

// Group is structure of Group.
type Group struct {
	ID       uuid.UUID
	SchoolID uuid.UUID
	Name     string
	GradeID  uuid.UUID

	ClassTeacherID *uuid.UUID
	ClassTeacher   *Teacher

	ClassPresidentID *uuid.UUID
	ClassPresident   *Student

	DeputyClassPresidentID *uuid.UUID
	DeputyClassPresident   *Student

	Lessons Lessons

	Grade Grade

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewGroup creates a new Group domain.
func NewGroup(
	schoolID uuid.UUID,
	name string,
	gradeID uuid.UUID,
	nowFunc func() time.Time,
) Group {
	now := nowFunc()

	return Group{
		ID:       uuid.New(),
		SchoolID: schoolID,
		Name:     name,
		GradeID:  gradeID,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// HasClassTeacher checks whether class teacher is assigned.
func (g *Group) HasClassTeacher() bool {
	return g.ClassTeacherID != nil
}

// HasClassPresident checks whether class president is assigned.
func (g *Group) HasClassPresident() bool {
	return g.ClassPresidentID != nil
}

// SetGrade sets grade to group.
func (g *Group) SetGrade(grade Grade) {
	g.Grade = grade
}

// SetClassTeacher sets class teacher to group.
func (g *Group) SetClassTeacher(teacher Teacher) {
	g.ClassTeacher = &teacher
}

// SetClassPresident sets class president to group.
func (g *Group) SetClassPresident(student Student) {
	g.ClassPresident = &student
}

// Update update group domain.
func (g *Group) Update(
	name string,
	gradeID uuid.UUID,

	classTeacherID *uuid.UUID,
	classPresidentID *uuid.UUID,
	deputyClassPresidentID *uuid.UUID,

	nowFunc func() time.Time,
) {
	now := nowFunc()

	g.Name = name
	g.GradeID = gradeID

	g.ClassTeacherID = classTeacherID
	g.ClassPresidentID = classPresidentID
	g.DeputyClassPresidentID = deputyClassPresidentID

	g.UpdatedAt = now
}

// GroupFilters is structure of Group filters.
type GroupFilters struct {
	ListFilter
}

// NewGroupFilters creates a new GroupFilters domain.
func NewGroupFilters() GroupFilters {
	return GroupFilters{}
}

// Groups is list of Group.
type Groups []Group

// GradeIDs returns a list of Group grades IDs.
func (g Groups) GradeIDs() []uuid.UUID {
	res := make([]uuid.UUID, 0, len(g))

	for _, group := range g {
		res = append(res, group.GradeID)
	}

	return res
}

// SetGrades sets grades to groups.
func (g Groups) SetGrades(grades Grades) {
	mapOfGradesByID := make(map[uuid.UUID]Grade, len(grades))
	for _, grade := range grades {
		mapOfGradesByID[grade.ID] = grade
	}

	for index := range g {
		g[index].SetGrade(mapOfGradesByID[g[index].GradeID])
	}
}
