package domain

import (
	"time"

	"github.com/google/uuid"
)

// Mark is student mark domain.
type Mark struct {
	ID          uuid.UUID
	LessonID    uuid.UUID
	StudentID   uuid.UUID
	Mark        string
	Description *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewMark creates a new mark domain.
func NewMark(
	lessonID uuid.UUID,
	studentID uuid.UUID,
	mark string,
	description *string,

	nowFunc func() time.Time,
) Mark {
	now := nowFunc()

	return Mark{
		ID:          uuid.New(),
		LessonID:    lessonID,
		StudentID:   studentID,
		Mark:        mark,
		Description: description,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Marks is slice of mark.
type Marks []Mark
