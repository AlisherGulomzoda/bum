package request

import (
	"time"

	"github.com/google/uuid"
)

// AssignWeekLessons is a request to add week lessons.
type AssignWeekLessons struct {
	GroupID    uuid.UUID    `json:"group_id" binding:"required,uuid"`
	WeekDate   time.Time    `json:"week_date" binding:"required"`
	LessonItem []LessonItem `json:"lessons"`
}

// LessonItem is lesson for adding.
type LessonItem struct {
	GroupSubjectID uuid.UUID  `json:"group_subject_id" binding:"required,uuid"`
	TeacherID      *uuid.UUID `json:"teacher_id" binding:"uuid"`
	AuditoriumID   uuid.UUID  `json:"auditorium_id" binding:"required,uuid"`
	StartTime      time.Time  `json:"start_time" binding:"required"`
	EndTime        time.Time  `json:"end_time" binding:"required"`
	Description    *string    `json:"description" binding:"required"`
}

// LessonsList request model for listing of lessons.
type LessonsList struct {
	ListFilter

	Period DateFilter

	SchoolID  *uuid.UUID `form:"school_id" binding:"omitempty,uuid"`
	TeacherID *uuid.UUID `form:"teacher_id" binding:"omitempty,uuid"`
	GroupID   *uuid.UUID `form:"group_id" binding:"omitempty,uuid"`
}
