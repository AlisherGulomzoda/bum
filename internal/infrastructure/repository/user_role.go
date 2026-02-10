package repository

import (
	"context"
	"fmt"
	"time"

	"bum-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	// UserRolesUserIDRoleSchoolIDOrganizationIDUniqueKey is combination of user_id, role and
	// school_id organization_id are unique key.
	UserRolesUserIDRoleSchoolIDOrganizationIDUniqueKey = "user_roles_user_id_role_school_id_organization_id_key"
	// UserRolesUserIDFKey is user id foreign key.
	UserRolesUserIDFKey = "user_roles_user_id_fkey"
	// UserRolesSchoolIDFKey is school id foreign key.
	UserRolesSchoolIDFKey = "user_roles_school_id_fkey"
	// UserRolesOrganizationIDFKey is school id foreign key.
	UserRolesOrganizationIDFKey = "user_roles_organization_id_fkey"
)

// UserRoleRow is a row containing user role.
type UserRoleRow struct {
	ID             uuid.UUID  `db:"id"`
	UserID         uuid.UUID  `db:"user_id"`
	Role           string     `db:"role"`
	SchoolID       *uuid.UUID `db:"school_id"`
	OrganizationID *uuid.UUID `db:"organization_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts an object into domain model.
func (e UserRoleRow) toDomain() domain.UserRole {
	return domain.UserRole{
		ID:             e.ID,
		UserID:         e.UserID,
		Role:           domain.Role(e.Role),
		SchoolID:       e.SchoolID,
		OrganizationID: e.OrganizationID,

		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

// UserRoleRows are a collection of UserRoleRow.
type UserRoleRows []UserRoleRow

func (u UserRoleRows) toDomain() domain.UserRoles {
	list := make(domain.UserRoles, 0, len(u))

	for index := range u {
		list = append(list, u[index].toDomain())
	}

	return list
}

// AddUserRoleTx creates a new user role within a transaction session.
func (u *User) AddUserRoleTx(ctx context.Context, d domain.UserRole, withinOnConflict bool) error {
	var (
		insertUserRoleQuery = `
			INSERT INTO user_roles
				( id,  user_id,  role,  school_id, organization_id )
			VALUES 
				(:id, :user_id, :role, :school_id,:organization_id )`

		args = map[string]any{
			"id":              d.ID,
			"user_id":         d.UserID,
			"role":            d.Role,
			"school_id":       d.SchoolID,
			"organization_id": d.OrganizationID,
		}
	)

	if withinOnConflict {
		insertUserRoleQuery += `
			ON CONFLICT ON CONSTRAINT user_roles_user_id_role_school_id_organization_id_key DO NOTHING`
	}

	_, err := u.session(ctx).NamedExecContext(ctx, insertUserRoleQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert user role: %w", err))
	}

	return nil
}

// UserRolesByIDTx get user roles by id.
func (u *User) UserRolesByIDTx(ctx context.Context, userID uuid.UUID) (domain.UserRoles, error) {
	var (
		getUserQuery = `
			SELECT 
				id, user_id, role, school_id, organization_id, created_at, updated_at, deleted_at
			FROM 
				user_roles
			WHERE 
				user_id = ? AND
				deleted_at IS NULL`

		row UserRoleRows
	)

	err := u.session(ctx).SelectContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getUserQuery), userID)
	if err != nil {
		return domain.UserRoles{}, handleError(fmt.Errorf("failed to select user roles by user_id: %w", err))
	}

	return row.toDomain(), nil
}

// UserRolesByIDsTx get user roles by ids.
func (u *User) UserRolesByIDsTx(ctx context.Context, userIDs []uuid.UUID) (domain.UserRoles, error) {
	var (
		getUserQuery = `
			SELECT 
				id, user_id, role, school_id, organization_id, created_at, updated_at, deleted_at
			FROM 
				user_roles
			WHERE 
				user_id IN (?) AND
				deleted_at IS NULL`

		row UserRoleRows
	)

	if len(userIDs) == 0 {
		return domain.UserRoles{}, nil
	}

	q, args, err := sqlx.In(getUserQuery, userIDs)
	if err != nil {
		return domain.UserRoles{}, handleError(fmt.Errorf("failed to select users roles by ids: %w", err))
	}

	err = u.session(ctx).SelectContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, q), args...)
	if err != nil {
		return domain.UserRoles{}, handleError(fmt.Errorf("failed to select user roles by user_id: %w", err))
	}

	return row.toDomain(), nil
}
