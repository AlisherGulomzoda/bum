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

// EduOrganization is educational organization repository.
type EduOrganization struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewEduOrganization create a new educational organization repository instance.
func NewEduOrganization(db postgres.DB, session transaction.SessionDB) *EduOrganization {
	return &EduOrganization{
		db:      db,
		session: session.DB,
	}
}

const (
	// EduOrganizationsNameKey is educational organization name unique key.
	EduOrganizationsNameKey = "educational_organizations_name_key"
)

// EduOrganizationRow is a row containing educational organization.
type EduOrganizationRow struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Logo        *string   `db:"logo"`
	Description *string   `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts an object into domain model.
func (e EduOrganizationRow) toDomain() domain.EduOrganization {
	return domain.EduOrganization{
		ID:          e.ID,
		Name:        e.Name,
		Logo:        e.Logo,
		Description: e.Description,

		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

// EduOrganizationRows is list of EduOrganizationRow.
type EduOrganizationRows []EduOrganizationRow

// toDomain converts object to entity.
func (e EduOrganizationRows) toDomain() domain.EduOrganizations {
	list := make(domain.EduOrganizations, 0, len(e))

	for _, row := range e {
		list = append(list, row.toDomain())
	}

	return list
}

// toShortDomain converts an object into domain model.
func (e EduOrganizationRow) toShortDomain() domain.EduOrganizationShortInfo {
	return domain.EduOrganizationShortInfo{
		ID:   e.ID,
		Name: e.Name,
	}
}

// toShortDomain converts an object into domain model.
func (e EduOrganizationRows) toShortDomain() domain.EduOrganizationShortInfos {
	list := make(domain.EduOrganizationShortInfos, 0, len(e))

	for _, row := range e {
		list = append(list, row.toShortDomain())
	}

	return list
}

// CreateEduOrganizationTx creates a new educational organization.
func (s *EduOrganization) CreateEduOrganizationTx(ctx context.Context, o domain.EduOrganization) error {
	var (
		sqlQuery = `
		INSERT INTO educational_organizations 
	    	( id, name, logo, description, created_at, updated_at ) 
		VALUES 
	    	(:id,:name,:logo,:description,:created_at,:updated_at );`

		args = map[string]any{
			"id":          o.ID,
			"name":        o.Name,
			"logo":        o.Logo,
			"description": o.Description,

			"created_at": o.CreatedAt,
			"updated_at": o.UpdatedAt,
		}
	)

	_, err := s.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert educational organization: %w", err))
	}

	return nil
}

// UpdateEduOrganizationTx update educational organization by id.
func (s *EduOrganization) UpdateEduOrganizationTx(ctx context.Context, o domain.EduOrganization) error {
	var (
		sqlQuery = `
		UPDATE educational_organizations 
			SET 
			    name=:name, 
			    logo=:logo, 
			    description=:description,
			    grade_standard_id=:grade_standard_id, 
			    updated_at=:updated_at
		WHERE
		    id = :id`

		args = map[string]any{
			"id":          o.ID,
			"name":        o.Name,
			"logo":        o.Logo,
			"description": o.Description,

			"updated_at": o.UpdatedAt,
		}
	)

	_, err := s.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert educational organization: %w", err))
	}

	return nil
}

// EduOrganizationListTx returns a list of educational organizations.
func (s *EduOrganization) EduOrganizationListTx(
	ctx context.Context, filters domain.EduOrganizationFilters,
) (domain.EduOrganizations, error) {
	var (
		sqlQuery = fmt.Sprintf(`
			SELECT
				id, name, logo, description, created_at, updated_at, deleted_at 
			FROM
				educational_organizations
			WHERE
				deleted_at IS NULL
			ORDER BY 
				created_at %s
			LIMIT ? OFFSET ?`,
			filters.SortOrder,
		)

		rows EduOrganizationRows
	)

	err := s.session(ctx).SelectContext(ctx, &rows, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), filters.Limit, filters.Offset)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select a list of educational organizations: %w", err))
	}

	return rows.toDomain(), nil
}

// EduOrganizationByIDTx get educational organization by id.
func (s *EduOrganization) EduOrganizationByIDTx(ctx context.Context, id uuid.UUID) (domain.EduOrganization, error) {
	var (
		sqlQuery = `
			SELECT
				id,
				name,
				logo,
				description,
				
				created_at,
				updated_at,
				deleted_at
			FROM 
				educational_organizations
			WHERE 
				id = ? AND 
				deleted_at IS NULL`

		row EduOrganizationRow
	)

	err := s.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), id)
	if err != nil {
		return domain.EduOrganization{},
			handleError(fmt.Errorf("failed to select educational organization by id: %w", err))
	}

	return row.toDomain(), nil
}

// EduOrganizationShortByIDTx get educational organization short info by id.
func (s *EduOrganization) EduOrganizationShortByIDTx(
	ctx context.Context, id uuid.UUID,
) (domain.EduOrganizationShortInfo, error) {
	var (
		sqlQuery = `
			SELECT
				id,
				name,
				
				created_at,
				updated_at,
				deleted_at
			FROM 
				educational_organizations
			WHERE 
				id = ? AND 
				deleted_at IS NULL`

		row EduOrganizationRow
	)

	err := s.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), id)
	if err != nil {
		return domain.EduOrganizationShortInfo{}, handleError(
			fmt.Errorf("failed to select educational organization by id: %w", err),
		)
	}

	return row.toShortDomain(), nil
}

// EduOrganizationsByIDsTx get educational organizations by ids.
func (s *EduOrganization) EduOrganizationsByIDsTx(
	ctx context.Context, ids []uuid.UUID,
) (domain.EduOrganizations, error) {
	var (
		sqlQuery = `
			SELECT
				id,
				name, 
				logo,
				description,
				
				created_at,
				updated_at,
				deleted_at
			FROM
				educational_organizations
			WHERE 
				id IN (?) AND 
				deleted_at IS NULL`

		rows EduOrganizationRows
	)

	if len(ids) == 0 {
		return domain.EduOrganizations{}, nil
	}

	q, args, err := sqlx.In(sqlQuery, ids)
	if err != nil {
		return domain.EduOrganizations{},
			handleError(fmt.Errorf("failed to select educational organization by ids: %w", err))
	}

	err = s.session(ctx).SelectContext(ctx, &rows, sqlx.Rebind(sqlx.DOLLAR, q), args...)
	if err != nil {
		return domain.EduOrganizations{},
			handleError(fmt.Errorf("failed to select educational organization by ids: %w", err))
	}

	return rows.toDomain(), nil
}

// EduOrganizationsShortInfoByIDsTx get educational organizations short info by ids.
func (s *EduOrganization) EduOrganizationsShortInfoByIDsTx(
	ctx context.Context, ids []uuid.UUID,
) (domain.EduOrganizationShortInfos, error) {
	var (
		sqlQuery = `
			SELECT
				id,
				name, 
				
				created_at,
				updated_at,
				deleted_at
			FROM
				educational_organizations
			WHERE 
				id IN (?) AND 
				deleted_at IS NULL`

		rows EduOrganizationRows
	)

	if len(ids) == 0 {
		return domain.EduOrganizationShortInfos{}, nil
	}

	q, args, err := sqlx.In(sqlQuery, ids)
	if err != nil {
		return domain.EduOrganizationShortInfos{},
			handleError(fmt.Errorf("failed to select educational organization by ids: %w", err))
	}

	err = s.session(ctx).SelectContext(ctx, &rows, sqlx.Rebind(sqlx.DOLLAR, q), args...)
	if err != nil {
		return domain.EduOrganizationShortInfos{},
			handleError(fmt.Errorf("failed to select educational organization by ids: %w", err))
	}

	return rows.toShortDomain(), nil
}

// EduOrganizationCountTx return educational organization count.
func (s *EduOrganization) EduOrganizationCountTx(ctx context.Context) (int, error) {
	var (
		sqlQuery = `
			SELECT
				count(*)
			FROM
				educational_organizations
			WHERE
				deleted_at IS NULL
	`

		count int
	)

	err := s.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, sqlQuery)).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get educational organization count: %w", err))
	}

	return count, nil
}
