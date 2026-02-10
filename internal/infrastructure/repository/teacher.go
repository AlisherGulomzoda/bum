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

// Teacher is teacher repository.
type Teacher struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewTeacher creates a new teacher repository.
func NewTeacher(db postgres.DB, session transaction.SessionDB) *Teacher {
	return &Teacher{
		db:      db,
		session: session.DB,
	}
}

const (
	// TeachersUserIDFKey is a teachers user_id foreign key.
	TeachersUserIDFKey = "teachers_user_id_fkey"
	// TeachersSchoolIDFKey is a teachers school_id foreign key.
	TeachersSchoolIDFKey = "teachers_school_id_fkey"
)

// TeacherRow is a teacher row.
type TeacherRow struct {
	ID       uuid.UUID `db:"id"`
	RoleID   uuid.UUID `db:"role_id"`
	UserID   uuid.UUID `db:"user_id"`
	SchoolID uuid.UUID `db:"school_id"`
	Phone    *string   `db:"phone"`
	Email    *string   `db:"email"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts an object into a domain model.
func (t TeacherRow) toDomain() domain.Teacher {
	return domain.Teacher{
		ID:       t.ID,
		RoleID:   t.RoleID,
		UserID:   t.UserID,
		SchoolID: t.SchoolID,
		Phone:    t.Phone,
		Email:    t.Email,

		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: t.DeletedAt,
	}
}

// TeacherRows is list of TeacherRow.
type TeacherRows []TeacherRow

// toDomain converts to entity.
func (t TeacherRows) toDomain() domain.Teachers {
	list := make(domain.Teachers, 0, len(t))

	for index := range t {
		list = append(list, t[index].toDomain())
	}

	return list
}

// CreateTeacherTx creates a new teacher in a transaction.
func (t *Teacher) CreateTeacherTx(ctx context.Context, o domain.Teacher) error {
	var (
		sqlQuery = `
			INSERT INTO teachers
				( id, role_id, user_id, school_id, phone, email, created_at, updated_at)
			VALUES	
				(:id,:role_id,:user_id,:school_id,:phone,:email,:created_at,:updated_at)`

		args = map[string]any{
			"id":        o.ID,
			"role_id":   o.RoleID,
			"user_id":   o.UserID,
			"school_id": o.SchoolID,
			"phone":     o.Phone,
			"email":     o.Email,

			"created_at": o.CreatedAt,
			"updated_at": o.UpdatedAt,
		}
	)

	_, err := t.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert teacher: %w", err))
	}

	return nil
}

// TeacherByIDTx get teacher by id.
func (t *Teacher) TeacherByIDTx(ctx context.Context, id uuid.UUID) (domain.Teacher, error) {
	var row TeacherRow

	getTeacherQuery := `
	SELECT 
		id,
		role_id,
		user_id,
		school_id, 
		phone, 
		email, 
		
		created_at, 
		updated_at, 
		deleted_at
	FROM 
		teachers
	WHERE 
		id = ? AND
		deleted_at IS NULL`

	err := t.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getTeacherQuery), id)
	if err != nil {
		return domain.Teacher{}, handleError(fmt.Errorf("failed to get teacher by id: %w", err))
	}

	return row.toDomain(), nil
}

// TeachersByIDsTx get teachers by ids.
func (t *Teacher) TeachersByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.Teachers, error) {
	var rows TeacherRows

	getTeacherQuery := `
	SELECT 
		id,
		role_id,
		user_id,
		school_id, 
		phone, 
		email, 
		
		created_at, 
		updated_at, 
		deleted_at
	FROM 
		teachers
	WHERE 
		id IN (?) AND
		deleted_at IS NULL`

	if len(ids) == 0 {
		return nil, nil
	}

	sqlQuery, params, err := sqlx.In(getTeacherQuery, ids)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select teacher by ids: %w", err))
	}

	err = t.session(ctx).SelectContext(ctx, &rows, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to get teachers by ids: %w", err))
	}

	return rows.toDomain(), nil
}

// TeacherListTx returns list of teachers by filter from database.
func (t *Teacher) TeacherListTx(ctx context.Context, filters domain.TeacherListFilter) (domain.Teachers, error) {
	params, filtersQuery, anySlices := teacherListFilter(filters)

	sqlQuery := `	
			SELECT 
				teachers.id,
				teachers.role_id,
				teachers.user_id,
				teachers.school_id, 
				teachers.phone,
				teachers.email,

				teachers.created_at, 
				teachers.updated_at, 
				teachers.deleted_at
			FROM 
				teachers
			INNER JOIN schools
			    ON schools.id = teachers.school_id
			LEFT JOIN group_subjects
				ON group_subjects.teacher_id = teachers.id
		` + where(filtersQuery) + ` GROUP BY teachers.id `

	sqlQuery += fmt.Sprintf(
		` ORDER BY teachers.created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		teacherList TeacherRows
		err         error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select teacher list: %w", err))
		}
	}

	err = t.session(ctx).SelectContext(ctx, &teacherList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select teacher list: %w", err))
	}

	return teacherList.toDomain(), nil
}

// teacherListFilter returns query by teacher list filter.
func teacherListFilter(filters domain.TeacherListFilter) (params []any, filtersQuery []string, anySlices bool) {
	filtersQuery = append(filtersQuery, "teachers.deleted_at IS NULL")
	anySlices = false

	if filters.CreatedDate.DateFrom != nil {
		filtersQuery = append(filtersQuery, "teachers.created_at >= ?")
		params = append(params, filters.CreatedDate.DateFrom)
	}

	if filters.CreatedDate.DateTill != nil {
		filtersQuery = append(filtersQuery, "teachers.created_at < ?")
		// add 1 day to include on director select, were added on this day.
		params = append(params, filters.CreatedDate.DateTill.AddDate(0, 0, 1))
	}

	if len(filters.SchoolIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "teachers.school_id IN (?)")
		params = append(params, filters.SchoolIDs)
	}

	if len(filters.GroupIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "group_subjects.group_id IN (?)")

		params = append(params, filters.GroupIDs)
	}

	if len(filters.OrganizationIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "schools.organization_id IN (?)")

		params = append(params, filters.OrganizationIDs)
	}

	return params, filtersQuery, anySlices
}

// TeacherCountTx returns list count of teacher by filter from database.
func (t *Teacher) TeacherCountTx(ctx context.Context, filters domain.TeacherListFilter) (int, error) {
	params, filtersQuery, anySlices := teacherListFilter(filters)

	sqlQuery := `
		SELECT COUNT(*) FROM ( SELECT 
			 teachers.*
		FROM 
			teachers
		INNER JOIN schools
			    ON schools.id = teachers.school_id
			LEFT JOIN group_subjects
				ON group_subjects.teacher_id = teachers.id
	` + where(filtersQuery) + ` GROUP BY teachers.id ) t`

	var (
		count int
		err   error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return 0, handleError(fmt.Errorf("failed to select teacher list count: %w", err))
		}
	}

	err = t.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get teacher list count: %w", err))
	}

	return count, nil
}
