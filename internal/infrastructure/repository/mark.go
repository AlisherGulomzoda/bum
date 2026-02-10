package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"bum-service/internal/domain"
)

const (
	// MarksStudentIDFKey is mark student id foreign key.
	MarksStudentIDFKey = "marks_student_id_fkey"
	// MarksLessonIDFKey is lesson_id foreign key.
	MarksLessonIDFKey = "marks_lesson_id_fkey"
	// MarksLessonKey is mark lesson key.
	MarksLessonKey = "marks_lesson_key"
)

// MarkRow is mark row.
type MarkRow struct {
	ID          uuid.UUID `db:"id"`
	LessonID    uuid.UUID `db:"lesson_id"`
	StudentID   uuid.UUID `db:"student_id"`
	Mark        string    `db:"mark"`
	Description *string   `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// MarkRows is slice of MarkRow.
type MarkRows []MarkRow

func (m MarkRows) toDomain() domain.Marks {
	res := make([]domain.Mark, 0, len(m))

	for _, row := range m {
		res = append(res, row.toDomain())
	}

	return res
}

func (m MarkRow) toDomain() domain.Mark {
	return domain.Mark{
		ID:          m.ID,
		LessonID:    m.LessonID,
		StudentID:   m.StudentID,
		Mark:        m.Mark,
		Description: m.Description,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

// AddMark adds mark to student.
func (l *Lesson) AddMark(ctx context.Context, m domain.Mark) error {
	var (
		sqlQuery = `
			INSERT INTO marks
				( id, lesson_id, student_id, mark, description, created_at, updated_at) 
			VALUES
				(:id,:lesson_id,:student_id,:mark,:description,:created_at,:updated_at) 
`

		args = map[string]any{
			"id":          m.ID,
			"lesson_id":   m.LessonID,
			"student_id":  m.StudentID,
			"mark":        m.Mark,
			"description": m.Description,

			"created_at": m.CreatedAt,
			"updated_at": m.UpdatedAt,
		}
	)

	_, err := l.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert mark: %w", err))
	}

	return nil
}

// MarkByIDTx get a mark by id.
func (l *Lesson) MarkByIDTx(ctx context.Context, id uuid.UUID) (domain.Mark, error) {
	var mark MarkRow

	sqlQuery := `
	SELECT 
		id, lesson_id, student_id, mark, description, created_at, updated_at, deleted_at 
	FROM 
	    marks 
	WHERE 
	    deleted_at IS NULL AND 
	    id = ?;
	`

	err := l.session(ctx).GetContext(ctx, &mark, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), id)
	if err != nil {
		return domain.Mark{}, handleError(fmt.Errorf("failed to select mark by id: %w", err))
	}

	return mark.toDomain(), nil
}
