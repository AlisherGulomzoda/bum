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

// Director is director repository.
type Director struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewDirector create a new director repository instance.
func NewDirector(db postgres.DB, session transaction.SessionDB) *Director {
	return &Director{
		db:      db,
		session: session.DB,
	}
}

const (
	// DirectorsUserIDFKey is directors user_id foreign key.
	DirectorsUserIDFKey = "directors_user_id_fkey"
	// DirectorsSchoolIDFKey is directors school_id foreign key.
	DirectorsSchoolIDFKey = "directors_school_id_fkey"
)

// DirectorRow is a row containing director.
type DirectorRow struct {
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
func (d DirectorRow) toDomain() domain.Director {
	return domain.Director{
		ID:       d.ID,
		RoleID:   d.RoleID,
		UserID:   d.UserID,
		SchoolID: d.SchoolID,
		Phone:    d.Phone,
		Email:    d.Email,

		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,
	}
}

// DirectorRows is list of DirectorRow.
type DirectorRows []DirectorRow

// toDomain converts to entity.
func (h DirectorRows) toDomain() domain.Directors {
	list := make(domain.Directors, 0, len(h))

	for index := range h {
		list = append(list, h[index].toDomain())
	}

	return list
}

// AddDirectorTx creates a new director within a transaction session.
func (d *Director) AddDirectorTx(ctx context.Context, o domain.Director) error {
	var (
		sqlQuery = `
			INSERT INTO directors
				(id, role_id, user_id, school_id, phone, email, created_at, updated_at) 
			VALUES	
				(:id,:role_id,:user_id,:school_id,:phone,:email,:created_at,:updated_at)`

		args = map[string]any{
			"id":         o.ID,
			"role_id":    o.RoleID,
			"user_id":    o.UserID,
			"school_id":  o.SchoolID,
			"phone":      o.Phone,
			"email":      o.Email,
			"created_at": o.CreatedAt,
			"updated_at": o.UpdatedAt,
		}
	)

	_, err := d.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert director: %w", err))
	}

	return nil
}

// whereAnd complements the selection for the query.
//
//nolint:unused // will be used.
func whereAnd(parameters []any) string {
	if len(parameters) == 0 {
		return ` WHERE`
	}

	return ` AND`
}

// DirectorListTx returns list of director by filter from database.
func (d *Director) DirectorListTx(ctx context.Context, filters domain.DirectorListFilter) (domain.Directors, error) {
	params, filtersQuery, anySlices := directorListFilter(filters)

	sqlQuery := `	
			SELECT 
				id, role_id, user_id, school_id, phone, email, created_at, updated_at, deleted_at
			FROM 
				directors
		` + where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		directorList = make(DirectorRows, 0)
		err          error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select director list: %w", err))
		}
	}

	err = d.session(ctx).SelectContext(ctx, &directorList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select director list: %w", err))
	}

	return directorList.toDomain(), nil
}

// directorListFilter returns query by director list filter.
func directorListFilter(filters domain.DirectorListFilter) (params []any, filtersQuery []string, anySlices bool) {
	filtersQuery = append(filtersQuery, "deleted_at IS NULL")

	if filters.CreatedDate.DateFrom != nil {
		filtersQuery = append(filtersQuery, "created_at >= ?")
		params = append(params, filters.CreatedDate.DateFrom)
	}

	if filters.CreatedDate.DateTill != nil {
		filtersQuery = append(filtersQuery, "created_at < ?")
		// add 1 day to include on director select, were added on this day.
		params = append(params, filters.CreatedDate.DateTill.AddDate(0, 0, 1))
	}

	if len(filters.SchoolIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "school_id IN (?)")

		params = append(params, filters.SchoolIDs)
	}

	return params, filtersQuery, anySlices
}

// DirectorCountTx returns count of director by filter from database.
func (d *Director) DirectorCountTx(ctx context.Context, filters domain.DirectorListFilter) (int, error) {
	params, filtersQuery, anySlices := directorListFilter(filters)

	q := `
		SELECT 
			COUNT(*) 
		FROM 
			directors
	` + where(filtersQuery)

	var (
		count int
		err   error
	)

	if anySlices {
		q, params, err = sqlx.In(q, params...)
		if err != nil {
			return count, handleError(fmt.Errorf("failed to get director list count : %w", err))
		}
	}

	err = d.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get director list count: %w", err))
	}

	return count, nil
}

// DirectorByIDTx get director by id.
func (d *Director) DirectorByIDTx(
	ctx context.Context, id uuid.UUID,
) (domain.Director, error) {
	var row DirectorRow

	getDirectorQuery := `
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
		directors
	WHERE 
		id = $1 AND
		deleted_at IS NULL`

	err := d.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getDirectorQuery), id)
	if err != nil {
		return domain.Director{}, handleError(fmt.Errorf("failed to select director by id: %w", err))
	}

	return row.toDomain(), nil
}
