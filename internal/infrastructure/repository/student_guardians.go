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
	// StudentGuardiansKey is student guardians unique key.
	StudentGuardiansKey = "student_guardians_key"
	// StudentGuardiansStudentIDFKey is student guardians student_id foreign key.
	StudentGuardiansStudentIDFKey = "student_guardians_student_id_fkey"
	// StudentGuardiansUserIDFKey is student guardians user_id foreign key.
	StudentGuardiansUserIDFKey = "student_guardians_user_id_fkey"
	// StudentGuardiansSchoolIDFKey is student guardians school_id foreign key.
	StudentGuardiansSchoolIDFKey = "student_guardians_school_id_fkey"
)

// StudentGuardianRow is a row containing student guardian.
type StudentGuardianRow struct {
	ID        uuid.UUID `db:"id"`
	StudentID uuid.UUID `db:"student_id"`
	UserID    uuid.UUID `db:"user_id"`
	SchoolID  uuid.UUID `db:"school_id"`
	Relation  string    `db:"relation"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// toDomain converts an object into domain model.
func (s StudentGuardianRow) toDomain() domain.StudentGuardian {
	return domain.StudentGuardian{
		ID:        s.ID,
		StudentID: s.StudentID,
		UserID:    s.UserID,
		SchoolID:  s.SchoolID,
		Relation:  domain.StudentGuardianRelation(s.Relation),

		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// StudentGuardianRows is list of StudentGuardianRow.
type StudentGuardianRows []StudentGuardianRow

// toDomain converts to entity.
func (s StudentGuardianRows) toDomain() domain.StudentGuardians {
	list := make(domain.StudentGuardians, 0, len(s))

	for index := range s {
		list = append(list, s[index].toDomain())
	}

	return list
}

// StudentGuardiansByStudentIDTx returns student guardians.
func (s *Student) StudentGuardiansByStudentIDTx(
	ctx context.Context,
	studentID uuid.UUID,
) (domain.StudentGuardians, error) {
	var (
		list = make(StudentGuardianRows, 0)
		err  error
	)

	query := `
		SELECT 
		    id, user_id, student_id, school_id, relation, created_at, updated_at
		FROM 
		    student_guardians
		WHERE 
		    student_id = ?`

	err = s.session(ctx).SelectContext(ctx, &list, sqlx.Rebind(sqlx.DOLLAR, query), studentID)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select student guardians by student id: %w", err))
	}

	return list.toDomain(), nil
}

// AddStudentGuardianTx adds student guardian.
func (s *Student) AddStudentGuardianTx(ctx context.Context, studentGuardian domain.StudentGuardian) error {
	var (
		sqlQuery = `
			INSERT INTO student_guardians
				(id, user_id, student_id, school_id, relation, created_at, updated_at)
			VALUES
				(:id,:user_id,:student_id,:school_id,:relation,:created_at,:updated_at)`

		args = map[string]any{
			"id":         studentGuardian.ID,
			"user_id":    studentGuardian.UserID,
			"student_id": studentGuardian.StudentID,
			"school_id":  studentGuardian.SchoolID,
			"relation":   studentGuardian.Relation,
			"created_at": studentGuardian.CreatedAt,
			"updated_at": studentGuardian.UpdatedAt,
		}
	)

	_, err := s.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert student guardian: %w", err))
	}

	return nil
}

// StudentGuardianByIDTx returns student guardian by id.
func (s *Student) StudentGuardianByIDTx(ctx context.Context, id uuid.UUID) (domain.StudentGuardian, error) {
	var row StudentGuardianRow

	getQuery := `
	SELECT 
		    id, user_id, student_id, school_id, relation, created_at, updated_at
	FROM 
		student_guardians
	WHERE 
		id = ?`

	err := s.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getQuery), id)
	if err != nil {
		return domain.StudentGuardian{}, handleError(
			fmt.Errorf("failed to select student guardian by id: %w", err))
	}

	return row.toDomain(), nil
}

// StudentGuardianByUserIDTx returns student guardian by user_id.
func (s *Student) StudentGuardianByUserIDTx(ctx context.Context, id uuid.UUID) (domain.StudentGuardians, error) {
	var rows StudentGuardianRows

	getQuery := `
	SELECT 
		    id, user_id, student_id, school_id, relation, created_at, updated_at
	FROM 
		student_guardians
	WHERE 
		user_id = ?`

	err := s.session(ctx).SelectContext(ctx, &rows, sqlx.Rebind(sqlx.DOLLAR, getQuery), id)
	if err != nil {
		return domain.StudentGuardians{}, handleError(
			fmt.Errorf("failed to select student guardian by user_id: %w", err))
	}

	return rows.toDomain(), nil
}

// StudentGuardianListTx returns list of student guardian by filter from database.
func (s *Student) StudentGuardianListTx(
	ctx context.Context,
	filters domain.StudentGuardianListFilter,
) (domain.StudentGuardians, error) {
	params, filtersQuery, anySlices := studentGuardianListFilter(filters)

	sqlQuery := `	
			SELECT 
				sg.id,
				sg.user_id,
				sg.student_id,
				sg.school_id,
				sg.relation,
				sg.created_at, 
				sg.updated_at
			FROM 
				student_guardians sg
			INNER JOIN
				schools ON schools.id = sg.school_id
			INNER JOIN
				students ON students.id = sg.student_id
		` + where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY sg.created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		list = make(StudentGuardianRows, 0)
		err  error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select student guardian list: %w", err))
		}
	}

	err = s.session(ctx).SelectContext(ctx, &list, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select student guardian list: %w", err))
	}

	return list.toDomain(), nil
}

// studentListFilter returns query by student guardian list filter.
func studentGuardianListFilter(
	filters domain.StudentGuardianListFilter,
) (params []any, filtersQuery []string, anySlices bool) {
	if filters.CreatedDate.DateFrom != nil {
		filtersQuery = append(filtersQuery, "sg.created_at >= ?")
		params = append(params, filters.CreatedDate.DateFrom)
	}

	if filters.CreatedDate.DateTill != nil {
		filtersQuery = append(filtersQuery, "sg.created_at < ?")
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

		filtersQuery = append(filtersQuery, "schools.id IN (?)")

		params = append(params, filters.SchoolIDs)
	}

	if len(filters.OrganizationIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "schools.organization_id IN (?)")

		params = append(params, filters.OrganizationIDs)
	}

	return params, filtersQuery, anySlices
}

// StudentGuardianListCountTx returns count of student guardian by filter from database.
func (s *Student) StudentGuardianListCountTx(
	ctx context.Context,
	filters domain.StudentGuardianListFilter,
) (int, error) {
	params, filtersQuery, anySlices := studentGuardianListFilter(filters)

	q := `
			SELECT 
				COUNT(*)
			FROM 
				student_guardians sg
			INNER JOIN
				schools ON schools.id = sg.school_id
			INNER JOIN
				students ON students.id = sg.student_id
	` + where(filtersQuery)

	var (
		count int
		err   error
	)

	if anySlices {
		q, params, err = sqlx.In(q, params...)
		if err != nil {
			return count, handleError(fmt.Errorf("failed to get student guardian list count: %w", err))
		}
	}

	err = s.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get student guardian list count: %w", err))
	}

	return count, nil
}
