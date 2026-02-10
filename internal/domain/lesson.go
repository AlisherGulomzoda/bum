package domain

import (
	"time"

	"github.com/google/uuid"
)

// Lesson is struct of lesson.
type Lesson struct {
	ID             uuid.UUID
	SchoolID       uuid.UUID
	GroupSubjectID uuid.UUID
	TeacherID      *uuid.UUID
	AuditoriumID   uuid.UUID
	StartTime      time.Time
	EndTime        time.Time
	Description    *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewLesson creates a new lesson domain.
func NewLesson(
	schoolID uuid.UUID,
	groupSubjectID uuid.UUID,
	teacherID *uuid.UUID,
	groupTeacherID *uuid.UUID,
	auditoriumID uuid.UUID,
	startTime time.Time,
	endTime time.Time,
	description *string,

	nowFunc func() time.Time,
) Lesson {
	now := nowFunc()

	// if there is no need for replacement then we set default teacher from group subject.
	if teacherID == nil {
		teacherID = groupTeacherID
	}

	return Lesson{
		ID:             uuid.New(),
		SchoolID:       schoolID,
		GroupSubjectID: groupSubjectID,
		TeacherID:      teacherID,
		AuditoriumID:   auditoriumID,
		StartTime:      startTime,
		EndTime:        endTime,
		Description:    description,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Lessons are list of lessons.
type Lessons []Lesson

// LessonsListFilter filter for the list of Lessons.
type LessonsListFilter struct {
	Period DateFilter
	ListFilter
	SchoolID  *uuid.UUID
	TeacherID *uuid.UUID
	GroupID   *uuid.UUID
}

// NewLessonsListFilter creates a new LessonsListFilter domain.
func NewLessonsListFilter(
	period DateFilter,
	list ListFilter,
	schoolID *uuid.UUID,
	teacherID *uuid.UUID,
	groupID *uuid.UUID,
) LessonsListFilter {
	return LessonsListFilter{
		Period:     period,
		ListFilter: list,
		SchoolID:   schoolID,
		TeacherID:  teacherID,
		GroupID:    groupID,
	}
}
