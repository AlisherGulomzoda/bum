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

// Headmaster is headmaster repository.
type Headmaster struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewHeadmaster create a new headmaster repository instance.
func NewHeadmaster(db postgres.DB, session transaction.SessionDB) *Headmaster {
	return &Headmaster{
		db:      db,
		session: session.DB,
	}
}

const (
	// HeadmastersUserIDFKey is a headmasters user_id foreign key.
	HeadmastersUserIDFKey = "headmasters_user_id_fkey"
	// HeadmastersSchoolIDFKey is headmasters school_id foreign key.
	HeadmastersSchoolIDFKey = "headmasters_school_id_fkey"
)

// HeadmasterRow is a row containing headmaster.
type HeadmasterRow struct {
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
func (h HeadmasterRow) toDomain() domain.Headmaster {
	return domain.Headmaster{
		ID:       h.ID,
		RoleID:   h.RoleID,
		UserID:   h.UserID,
		SchoolID: h.SchoolID,
		Phone:    h.Phone,
		Email:    h.Email,

		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
		DeletedAt: h.DeletedAt,
	}
}

// HeadmasterRows is list of HeadmasterRow.
type HeadmasterRows []HeadmasterRow

// toDomain converts to entity.
func (h HeadmasterRows) toDomain() domain.Headmasters {
	list := make(domain.Headmasters, 0, len(h))

	for index := range h {
		list = append(list, h[index].toDomain())
	}

	return list
}

// AddHeadmasterTx creates a new headmaster within a transaction session.
func (h *Headmaster) AddHeadmasterTx(ctx context.Context, o domain.Headmaster) error {
	var (
		sqlQuery = `
			INSERT INTO headmasters
				( id, role_id, user_id, school_id, phone, email, created_at, updated_at )
			VALUES
				(:id,:role_id,:user_id,:school_id,:phone,:email,:created_at,:updated_at );`

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

	_, err := h.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert headmaster: %w", err))
	}

	return nil
}

// HeadmasterByIDTx get headmaster by id.
func (h *Headmaster) HeadmasterByIDTx(ctx context.Context, id uuid.UUID) (domain.Headmaster, error) {
	var (
		getHeadmasterQuery = `
			SELECT 
				id, role_id, user_id, school_id, phone, email, created_at, updated_at, deleted_at
			FROM 
				headmasters
			WHERE 
				id = ? AND 
				deleted_at IS NULL`

		row HeadmasterRow
	)

	err := h.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getHeadmasterQuery), id)
	if err != nil {
		return domain.Headmaster{}, handleError(fmt.Errorf("failed to get headmaster by id: %w", err))
	}

	return row.toDomain(), nil
}

// HeadmasterListTx returns list of headmaster by filter from database.
func (h *Headmaster) HeadmasterListTx(
	ctx context.Context,
	filters domain.HeadmasterListFilter,
) (domain.Headmasters, error) {
	params, filtersQuery, anySlices := headmasterListFilter(filters)

	sqlQuery := `  SELECT 
				id, role_id, user_id, school_id, phone, email, created_at, updated_at, deleted_at
			FROM 
				headmasters
		` + where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		headmasterList HeadmasterRows
		err            error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select headmaster list: %w", err))
		}
	}

	err = h.session(ctx).SelectContext(ctx, &headmasterList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select headmaster list: %w", err))
	}

	return headmasterList.toDomain(), nil
}

// headmasterListFilter returns query by headmaster list filter.
func headmasterListFilter(filters domain.HeadmasterListFilter) (params []any, filtersQuery []string, anySlices bool) {
	filtersQuery = append(filtersQuery, "deleted_at IS NULL")
	anySlices = false

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

// HeadmasterCountTx returns count of headmaster by filter from database.
func (h *Headmaster) HeadmasterCountTx(ctx context.Context, filters domain.HeadmasterListFilter) (int, error) {
	params, filtersQuery, anySlices := headmasterListFilter(filters)

	sqlQuery := `
		SELECT 
			COUNT(*) 
		FROM 
			headmasters
	` + where(filtersQuery)

	var (
		count int
		err   error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return 0, handleError(fmt.Errorf("failed to get headmaster list count: %w", err))
		}
	}

	err = h.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get headmaster list count: %w", err))
	}

	return count, nil
}
