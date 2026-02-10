package domain

import (
	"time"

	"github.com/google/uuid"
)

// GroupSubject is a group subject domain.
type GroupSubject struct {
	ID              uuid.UUID
	SchoolSubjectID uuid.UUID
	GroupID         uuid.UUID
	TeacherID       *uuid.UUID
	Count           *int16

	SchoolSubject SchoolSubject

	Teacher *Teacher

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewGroupSubject creates a new group subject.
func NewGroupSubject(
	schoolSubjectID,
	groupID uuid.UUID,
	teacherID *uuid.UUID,
	count *int16,
	nowFunc func() time.Time,
) GroupSubject {
	now := nowFunc()

	return GroupSubject{
		ID:              uuid.New(),
		SchoolSubjectID: schoolSubjectID,
		GroupID:         groupID,
		TeacherID:       teacherID,
		Count:           count,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// HasTeacher checks whether teacher is assigned to subject.
func (g *GroupSubject) HasTeacher() bool {
	return g.TeacherID != nil
}

// SetTeacher sets teacher to group subject.
func (g *GroupSubject) SetTeacher(teacher Teacher) {
	g.Teacher = &teacher
}

// GroupSubjects is a group subjects domain model.
type GroupSubjects []GroupSubject

// MapByID returns map of group subject by id.
func (g GroupSubjects) MapByID() map[uuid.UUID]GroupSubject {
	m := make(map[uuid.UUID]GroupSubject, len(g))

	for _, gs := range g {
		m[gs.ID] = gs
	}

	return m
}

// SchoolSubjectIDs returns list of school subject ids.
func (g GroupSubjects) SchoolSubjectIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(g))

	for _, subject := range g {
		list = append(list, subject.SchoolSubjectID)
	}

	return list
}

// TeacherIDs returns list of teacher ids.
func (g GroupSubjects) TeacherIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(g))

	for _, subject := range g {
		if subject.TeacherID != nil {
			list = append(list, *subject.TeacherID)
		}
	}

	return list
}

// SetSchoolSubjects sets school subjects to group subjects.
func (g GroupSubjects) SetSchoolSubjects(schoolSubjects SchoolSubjects) {
	mapOfSchoolSubjects := make(map[uuid.UUID]SchoolSubject, len(schoolSubjects))
	for index := range schoolSubjects {
		mapOfSchoolSubjects[schoolSubjects[index].ID] = schoolSubjects[index]
	}

	for index := 0; index < len(g); index++ {
		g[index].SchoolSubject = mapOfSchoolSubjects[g[index].SchoolSubjectID]
	}
}

// SetTeachers sets teachers to subjects.
func (g GroupSubjects) SetTeachers(teachers Teachers) {
	mapOfTeachers := make(map[uuid.UUID]*Teacher, len(teachers))
	for index := range teachers {
		mapOfTeachers[teachers[index].ID] = &teachers[index]
	}

	for index := 0; index < len(g); index++ {
		teacherID := g[index].TeacherID
		if teacherID != nil {
			g[index].Teacher = mapOfTeachers[*teacherID]
		}
	}
}
