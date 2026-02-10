//nolint:dupl // it's ok
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

// Student is student repository.
type Student struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewStudent create a new student repository instance.
func NewStudent(db postgres.DB, session transaction.SessionDB) *Student {
	return &Student{
		db:      db,
		session: session.DB,
	}
}

const (
	// StudentsGroupIDFKey is a student group_id foreign key.
	StudentsGroupIDFKey = "students_group_id_fkey"
	// StudentsUserIDFKey is a student user_id foreign key.
	StudentsUserIDFKey = "students_user_id_fkey"
)

// StudentRow is a row containing student.
type StudentRow struct {
	ID       uuid.UUID `db:"id"`
	RoleID   uuid.UUID `db:"role_id"`
	UserID   uuid.UUID `db:"user_id"`
	GroupID  uuid.UUID `db:"group_id"`
	SchoolID uuid.UUID `db:"school_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts an object into a domain model.
func (s StudentRow) toDomain() domain.Student {
	return domain.Student{
		ID:       s.ID,
		RoleID:   s.RoleID,
		UserID:   s.UserID,
		GroupID:  s.GroupID,
		SchoolID: s.SchoolID,

		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		DeletedAt: s.DeletedAt,
	}
}

// StudentRows is list of StudentRow.
type StudentRows []StudentRow

// toDomain converts to entity.
func (s StudentRows) toDomain() domain.Students {
	list := make(domain.Students, 0, len(s))

	for index := range s {
		list = append(list, s[index].toDomain())
	}

	return list
}

// AddStudentTx creates a new student within a transaction session.
func (s *Student) AddStudentTx(ctx context.Context, o domain.Student) error {
	var (
		sqlQuery = `
			INSERT INTO students
				( id, role_id, user_id, group_id, created_at, updated_at )
			VALUES
				(:id,:role_id,:user_id,:group_id,:created_at,:updated_at );`

		args = map[string]any{
			"id":       o.ID,
			"role_id":  o.RoleID,
			"user_id":  o.UserID,
			"group_id": o.GroupID,

			"created_at": o.CreatedAt,
			"updated_at": o.UpdatedAt,
		}
	)

	_, err := s.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert student: %w", err))
	}

	return nil
}

// StudentsByIDsTx get students by ids.
func (s *Student) StudentsByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.Students, error) {
	var (
		getQuery = `
			SELECT 
				id,
				"role_id",
				user_id,
				group_id,
				
				created_at, 
				updated_at, 
				deleted_at
			FROM 
				students
			WHERE 
				id IN (?) AND
				deleted_at IS NULL`

		rows StudentRows
	)

	if len(ids) == 0 {
		return domain.Students{}, nil
	}

	q, args, err := sqlx.In(getQuery, ids)
	if err != nil {
		return domain.Students{}, handleError(fmt.Errorf("failed to select students by ids: %w", err))
	}

	err = s.session(ctx).SelectContext(ctx, &rows, sqlx.Rebind(sqlx.DOLLAR, q), args...)
	if err != nil {
		return domain.Students{}, handleError(fmt.Errorf("failed to select students by ids: %w", err))
	}

	return rows.toDomain(), nil
}

// StudentByIDTx get student by id.
func (s *Student) StudentByIDTx(ctx context.Context, id uuid.UUID) (domain.Student, error) {
	var row StudentRow

	getQuery := `
	SELECT 
		students.id,
		students.role_id,
		students.user_id,
		students.group_id,
		groups.school_id,
		
		students.created_at, 
		students.updated_at, 
		students.deleted_at
	FROM 
		students
	INNER JOIN 
		    groups ON students.group_id = groups.id
	WHERE 
		students.id = $1 AND
		students.deleted_at IS NULL`

	err := s.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getQuery), id)
	if err != nil {
		return domain.Student{}, handleError(fmt.Errorf("failed to select student by id: %w", err))
	}

	return row.toDomain(), nil
}

// StudentsByUserIDTx get student by user_id.
func (s *Student) StudentsByUserIDTx(ctx context.Context, userID uuid.UUID) (domain.Students, error) {
	var rows StudentRows

	getQuery := `
	SELECT 
		students.id,
		students.role_id,
		students.user_id,
		students.group_id,
		groups.school_id,
		
		students.created_at, 
		students.updated_at, 
		students.deleted_at
	FROM 
		students
	INNER JOIN 
		    groups ON students.group_id = groups.id
	WHERE 
		students.user_id = $1 AND
		students.deleted_at IS NULL`

	err := s.session(ctx).SelectContext(ctx, &rows, sqlx.Rebind(sqlx.DOLLAR, getQuery), userID)
	if err != nil {
		return domain.Students{}, handleError(fmt.Errorf("failed to select student by user_id: %w", err))
	}

	return rows.toDomain(), nil
}

// StudentListTx returns list of students by filter from database.
func (s *Student) StudentListTx(ctx context.Context, filters domain.StudentListFilter) (domain.Students, error) {
	params, filtersQuery, anySlices := studentListFilter(filters)

	sqlQuery := `	
			SELECT 
				students.id,
				students.role_id,
				students.user_id,
				students.group_id,
				groups.school_id,

				students.created_at, 
				students.updated_at, 
				students.deleted_at
			FROM 
				students
			INNER JOIN
				groups ON groups.id = students.group_id
			INNER JOIN
				schools ON schools.id = groups.school_id
		` + where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY students.created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		studentList = make(StudentRows, 0)
		err         error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select student list: %w", err))
		}
	}

	err = s.session(ctx).SelectContext(ctx, &studentList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select student list: %w", err))
	}

	return studentList.toDomain(), nil
}

// studentListFilter returns query by student list filter.
func studentListFilter(filters domain.StudentListFilter) (params []any, filtersQuery []string, anySlices bool) {
	filtersQuery = append(filtersQuery, "students.deleted_at IS NULL")

	if filters.CreatedDate.DateFrom != nil {
		filtersQuery = append(filtersQuery, "students.created_at >= ?")
		params = append(params, filters.CreatedDate.DateFrom)
	}

	if filters.CreatedDate.DateTill != nil {
		filtersQuery = append(filtersQuery, "students.created_at < ?")
		// add 1 day to include on director select, were added on this day.
		params = append(params, filters.CreatedDate.DateTill.AddDate(0, 0, 1))
	}

	if len(filters.GroupIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "students.group_id IN (?)")

		params = append(params, filters.GroupIDs)
	}

	if len(filters.SchoolIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "groups.school_id IN (?)")

		params = append(params, filters.SchoolIDs)
	}

	if len(filters.OrganizationIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "schools.organization_id IN (?)")

		params = append(params, filters.OrganizationIDs)
	}

	return params, filtersQuery, anySlices
}

// StudentCountTx returns count of student by filter from database.
func (s *Student) StudentCountTx(ctx context.Context, filters domain.StudentListFilter) (int, error) {
	params, filtersQuery, anySlices := studentListFilter(filters)

	q := `
			SELECT 
				COUNT(*)
			FROM 
				students
			INNER JOIN
				groups ON groups.id = students.group_id
			INNER JOIN
				schools ON schools.id = groups.school_id
	` + where(filtersQuery)

	var (
		count int
		err   error
	)

	if anySlices {
		q, params, err = sqlx.In(q, params...)
		if err != nil {
			return count, handleError(fmt.Errorf("failed to get student list count : %w", err))
		}
	}

	err = s.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get student list count: %w", err))
	}

	return count, nil
}
