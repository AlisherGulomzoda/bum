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

// Owner is Owner repository.
type Owner struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewOwner create a new owner repository instance.
func NewOwner(db postgres.DB, session transaction.SessionDB) *Owner {
	return &Owner{
		db:      db,
		session: session.DB,
	}
}

const (
	// OwnersUserIDFKey is owners user_id foreign key.
	OwnersUserIDFKey = "owners_user_id_fkey"
	// OwnersOrganizationIDFKey is owners organization_id foreign key.
	OwnersOrganizationIDFKey = "owners_organization_id_fkey"
)

// OwnerRow is a row containing owner.
type OwnerRow struct {
	ID             uuid.UUID `db:"id"`
	RoleID         uuid.UUID `db:"role_id"`
	UserID         uuid.UUID `db:"user_id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Phone          *string   `db:"phone"`
	Email          *string   `db:"email"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts an object into a domain model.
func (o OwnerRow) toDomain() domain.Owner {
	return domain.Owner{
		ID:             o.ID,
		RoleID:         o.RoleID,
		UserID:         o.UserID,
		OrganizationID: o.OrganizationID,
		Phone:          o.Phone,
		Email:          o.Email,

		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		DeletedAt: o.DeletedAt,
	}
}

// OwnerRows is list of OwnerRow.
type OwnerRows []OwnerRow

// toDomain converts to entity.
func (o OwnerRows) toDomain() domain.Owners {
	list := make(domain.Owners, 0, len(o))

	for index := range o {
		list = append(list, o[index].toDomain())
	}

	return list
}

// AddOwnerTx creates a new owner within a transaction session.
func (o *Owner) AddOwnerTx(ctx context.Context, owner domain.Owner) error {
	var (
		sqlQuery = `
			INSERT INTO owners
				( id, role_id, user_id, organization_id, phone, email, created_at, updated_at )
			VALUES
				(:id,:role_id,:user_id,:organization_id,:phone,:email,:created_at,:updated_at );`

		args = map[string]any{
			"id":              owner.ID,
			"role_id":         owner.RoleID,
			"user_id":         owner.UserID,
			"organization_id": owner.OrganizationID,
			"phone":           owner.Phone,
			"email":           owner.Email,

			"created_at": owner.CreatedAt,
			"updated_at": owner.UpdatedAt,
		}
	)

	_, err := o.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert headmaster: %w", err))
	}

	return nil
}

// OwnerByIDTx get owner by id.
func (o *Owner) OwnerByIDTx(ctx context.Context, id uuid.UUID) (domain.Owner, error) {
	var (
		sqlQuery = `
			SELECT 
				id, role_id, user_id, organization_id, phone, email, created_at, updated_at, deleted_at
			FROM 
				owners
			WHERE 
				id = ? AND 
				deleted_at IS NULL`

		row OwnerRow
	)

	err := o.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), id)
	if err != nil {
		return domain.Owner{}, handleError(fmt.Errorf("failed to get owner by id: %w", err))
	}

	return row.toDomain(), nil
}

// OwnerByUserIDAndSchoolIDTx get owner by user_id and school_id.
func (o *Owner) OwnerByUserIDAndSchoolIDTx(ctx context.Context, schoolID, userID uuid.UUID) (domain.Owner, error) {
	var (
		sqlQuery = `
			SELECT 
				id, role_id, user_id, organization_id, phone, email, created_at, updated_at, deleted_at
			FROM 
				owners
			WHERE 
				user_id = ? AND
				school_id = ? AND
				deleted_at IS NULL`

		row OwnerRow
	)

	err := o.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), schoolID, userID)
	if err != nil {
		return domain.Owner{}, handleError(fmt.Errorf("failed to get owner by user_id and school_id: %w", err))
	}

	return row.toDomain(), nil
}

// OwnerListTx returns list of owner by filter from database.
func (o *Owner) OwnerListTx(ctx context.Context, filters domain.OwnerListFilter) (domain.Owners, error) {
	params, filtersQuery, anySlices := ownerListFilter(filters)

	sqlQuery := `	
			SELECT 
				id, role_id, user_id, organization_id, phone, email, created_at, updated_at, deleted_at
			FROM 
				owners
		` + where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		ownersList = make(OwnerRows, 0)
		err        error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select owner list: %w", err))
		}
	}

	err = o.session(ctx).SelectContext(ctx, &ownersList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select owner list: %w", err))
	}

	return ownersList.toDomain(), nil
}

// ownerListFilter returns query by owner list filter.
func ownerListFilter(filters domain.OwnerListFilter) (params []any, filtersQuery []string, anySlices bool) {
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

	if len(filters.OrganizationIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "organization_id IN (?)")

		params = append(params, filters.OrganizationIDs)
	}

	return params, filtersQuery, anySlices
}

// OwnerCountTx returns count of owners by filter from database.
func (o *Owner) OwnerCountTx(ctx context.Context, filters domain.OwnerListFilter) (int, error) {
	params, filtersQuery, anySlices := ownerListFilter(filters)

	q := `
		SELECT 
			COUNT(*) 
		FROM 
			owners
	` + where(filtersQuery)

	var (
		count int
		err   error
	)

	if anySlices {
		q, params, err = sqlx.In(q, params...)
		if err != nil {
			return count, handleError(fmt.Errorf("failed to get owners list count : %w", err))
		}
	}

	err = o.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get owners list count: %w", err))
	}

	return count, nil
}
