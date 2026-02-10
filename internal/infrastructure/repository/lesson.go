package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"bum-service/internal/domain"
	"bum-service/pkg/postgres"
	"bum-service/pkg/transaction"
)

// Lesson is lesson repository.
type Lesson struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewLesson creates a new lesson repository.
func NewLesson(db postgres.DB, session transaction.SessionDB) *Lesson {
	return &Lesson{
		db:      db,
		session: session.DB,
	}
}

const (
	// LessonsSchoolIDFKey is lesson school id foreign key.
	LessonsSchoolIDFKey = "lessons_school_id_fkey"
	// LessonsGroupSubjectIDFKey is lesson group subject id foreign key.
	LessonsGroupSubjectIDFKey = "lessons_group_subject_id_fkey"
	// LessonsTeacherIDFKey is teacher id foreign key.
	LessonsTeacherIDFKey = "lessons_teacher_id_fkey"
	// LessonsAuditoriumIDFKey is auditorium id foreign key.
	LessonsAuditoriumIDFKey = "lessons_auditorium_id_fkey"
)

// LessonRow is row containing lesson.
type LessonRow struct {
	ID             uuid.UUID  `db:"id"`
	SchoolID       uuid.UUID  `db:"school_id"`
	GroupSubjectID uuid.UUID  `db:"group_subject_id"`
	TeacherID      *uuid.UUID `db:"teacher_id"`
	AuditoriumID   uuid.UUID  `db:"auditorium_id"`
	StartTime      time.Time  `db:"start_time"`
	EndTime        time.Time  `db:"end_time"`
	Description    *string    `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// LessonRows is list of lesson rows.
type LessonRows []LessonRow

func (l LessonRow) toDomain() domain.Lesson {
	return domain.Lesson{
		ID:             l.ID,
		SchoolID:       l.SchoolID,
		GroupSubjectID: l.GroupSubjectID,
		TeacherID:      l.TeacherID,
		AuditoriumID:   l.AuditoriumID,
		StartTime:      l.StartTime,
		EndTime:        l.EndTime,
		Description:    l.Description,
		CreatedAt:      l.CreatedAt,
		UpdatedAt:      l.UpdatedAt,
		DeletedAt:      l.DeletedAt,
	}
}

func (l LessonRows) toDomain() domain.Lessons {
	lessons := make(domain.Lessons, 0, len(l))

	for _, lesson := range l {
		lessons = append(lessons, lesson.toDomain())
	}

	return lessons
}

// AssignsLessons adds lessons for a week.
func (l *Lesson) AssignsLessons(
	ctx context.Context, groupID uuid.UUID, firstDayOfWeek, firstDayOfNextWeek time.Time, lessons domain.Lessons,
) error {
	var (
		deleteQuery = `
			DELETE FROM 
			           lessons AS l
			       USING 
			           group_subjects AS gs
			WHERE 
			    gs.id = l.group_subject_id AND 
			    gs.group_id = :group_id AND
			    l.start_time >= :first_day_of_week AND 
			    l.end_time < :first_day_of_next_week
		`

		insertQuery = `
	INSERT INTO 
			lessons
		( 
			id, school_id, group_subject_id, teacher_id, auditorium_id, 
			start_time, end_time, description, created_at, updated_at
		) 
	VALUES 
		(
			:id,:school_id,:group_subject_id,:teacher_id,:auditorium_id,
			:start_time,:end_time,:description,:created_at,:updated_at
		)`
	)

	_, err := l.session(ctx).NamedExecContext(
		ctx, deleteQuery,
		map[string]any{
			"group_id":               groupID,
			"first_day_of_week":      firstDayOfWeek,
			"first_day_of_next_week": firstDayOfNextWeek,
		})
	if err != nil {
		return handleError(fmt.Errorf("failed to remove group subjects : %w", err))
	}

	listOfLessonsInsertRows := make([]map[string]any, 0, len(lessons))

	for _, lesson := range lessons {
		lessonsInsertRow := make(map[string]any)
		lessonsInsertRow["id"] = lesson.ID
		lessonsInsertRow["school_id"] = lesson.SchoolID
		lessonsInsertRow["group_subject_id"] = lesson.GroupSubjectID
		lessonsInsertRow["teacher_id"] = lesson.TeacherID
		lessonsInsertRow["auditorium_id"] = lesson.AuditoriumID
		lessonsInsertRow["start_time"] = lesson.StartTime
		lessonsInsertRow["end_time"] = lesson.EndTime
		lessonsInsertRow["description"] = lesson.Description
		lessonsInsertRow["created_at"] = lesson.CreatedAt
		lessonsInsertRow["updated_at"] = lesson.UpdatedAt

		listOfLessonsInsertRows = append(listOfLessonsInsertRows, lessonsInsertRow)
	}

	_, err = l.session(ctx).NamedExecContext(ctx, insertQuery, listOfLessonsInsertRows)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert lessons : %w", err))
	}

	return nil
}

// LessonsListTx returns list of lessons by filter from database.
func (l *Lesson) LessonsListTx(
	ctx context.Context,
	filters domain.LessonsListFilter,
) (domain.Lessons, error) {
	params, filtersQuery, anySlices := lessonsListFilter(filters)

	sqlQuery := `  
		SELECT 
			l.id, 
			l.school_id, 
			l.group_subject_id, 
			l.teacher_id, 
			l.auditorium_id, 
			l.start_time, 
			l.end_time, 
			l.description, 
			l.created_at, 
			l.updated_at
		FROM 
			lessons AS l
		INNER JOIN 
			group_subjects AS gs ON l.group_subject_id = gs.id
		` + where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY l.start_time %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		lessonsList LessonRows
		err         error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select headmaster list: %w", err))
		}
	}

	err = l.session(ctx).SelectContext(ctx, &lessonsList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select lessons list: %w", err))
	}

	return lessonsList.toDomain(), nil
}

// lessonsListFilter returns query by headmaster list filter.
func lessonsListFilter(filters domain.LessonsListFilter) (params []any, filtersQuery []string, anySlices bool) {
	filtersQuery = append(filtersQuery, "l.deleted_at IS NULL")
	anySlices = false

	if filters.Period.DateFrom != nil {
		filtersQuery = append(filtersQuery, "l.start_time >= ?")
		params = append(params, filters.Period.DateFrom)
	}

	if filters.Period.DateTill != nil {
		filtersQuery = append(filtersQuery, "l.end_time < ?")
		// add 1 day to include on director select, were added on this day.
		params = append(params, filters.Period.DateTill.AddDate(0, 0, 1))
	}

	if filters.SchoolID != nil {
		filtersQuery = append(filtersQuery, "l.school_id = ?")
		params = append(params, filters.SchoolID)
	}

	if filters.GroupID != nil {
		filtersQuery = append(filtersQuery, "gs.group_id = ?")
		params = append(params, filters.GroupID)
	}

	if filters.TeacherID != nil {
		filtersQuery = append(filtersQuery, "l.teacher_id = ?")
		params = append(params, filters.TeacherID, filters.TeacherID)
	}

	return params, filtersQuery, anySlices
}
