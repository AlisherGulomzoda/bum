package response

import (
	"time"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// Lesson is lessons response item.
type Lesson struct {
	ID             uuid.UUID  `json:"id"`
	GroupSubjectID uuid.UUID  `json:"group_subject_id"`
	TeacherID      *uuid.UUID `json:"teacher_id"`
	AuditoriumID   uuid.UUID  `json:"auditorium_id"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        time.Time  `json:"end_time"`
	Description    *string    `json:"description"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewLesson converts domain lesson into response lesson.
func NewLesson(lesson domain.Lesson) Lesson {
	return Lesson{
		ID:             lesson.ID,
		GroupSubjectID: lesson.GroupSubjectID,
		TeacherID:      lesson.TeacherID,
		AuditoriumID:   lesson.AuditoriumID,
		StartTime:      lesson.StartTime,
		EndTime:        lesson.EndTime,
		Description:    lesson.Description,

		CreatedAt: utils.RFC3339Time(lesson.CreatedAt),
		UpdatedAt: utils.RFC3339Time(lesson.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(lesson.DeletedAt),
	}
}

// NewLessons creates a new response of lesson item.
func NewLessons(list domain.Lessons) any {
	lessons := make([]Lesson, 0, len(list))
	for _, lesson := range list {
		lessons = append(lessons, NewLesson(lesson))
	}

	return lessons
}
