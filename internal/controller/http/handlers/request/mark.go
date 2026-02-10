package request

import "github.com/google/uuid"

// AddMark is mark request for adding mark.
type AddMark struct {
	LessonID    uuid.UUID `json:"lesson_id" binding:"required,uuid"`
	StudentID   uuid.UUID `json:"student_id" binding:"required,uuid"`
	Mark        string    `json:"mark" binding:"required"`
	Description *string   `json:"description"`
}
