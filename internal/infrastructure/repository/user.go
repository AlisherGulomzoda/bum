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

// User is users repository.
type User struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewUser create a new user repository instance.
func NewUser(db postgres.DB, session transaction.SessionDB) *User {
	return &User{
		db:      db,
		session: session.DB,
	}
}

const (
	// UsersEmailUniqueKey is users email unique key.
	UsersEmailUniqueKey = "users_email_key"
	// UsersPhoneUniqueKey is users phone unique key.
	UsersPhoneUniqueKey = "users_phone_key"
)

// UserRow is a row containing user.
type UserRow struct {
	ID         uuid.UUID  `db:"id"`
	FirstName  string     `db:"first_name"`
	LastName   string     `db:"last_name"`
	MiddleName *string    `db:"middle_name"`
	Gender     string     `db:"gender"`
	Phone      *string    `db:"phone"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

// toDomain converts an object into a domain model.
func (e UserRow) toDomain() domain.User {
	return domain.User{
		ID:         e.ID,
		FirstName:  e.FirstName,
		LastName:   e.LastName,
		MiddleName: e.MiddleName,
		Gender:     domain.Gender(e.Gender),
		Phone:      e.Phone,
		Email:      e.Email,

		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

// toDomain converts an object into domain model with password.
func (e UserRow) toDomainWithPassword() domain.User {
	// TODO поправить костыль
	domainUser := e.toDomain()
	domainUser.Password = e.Password

	return domainUser
}

// UserRows are a collection of user row.
type UserRows []UserRow

func (u UserRows) toDomain() domain.Users {
	list := make(domain.Users, 0, len(u))

	for index := range u {
		list = append(list, u[index].toDomain())
	}

	return list
}

// AddUserTx creates a new user within a transaction session.
func (u *User) AddUserTx(ctx context.Context, o domain.User) error {
	var (
		insertUserQuery = `
			INSERT INTO users
				( id, first_name, last_name, middle_name, gender, phone, email, password, created_at, updated_at ) 
			VALUES 
				(:id,:first_name,:last_name,:middle_name,:gender,:phone,:email,:password,:created_at,:updated_at );`

		args = map[string]any{
			"id":          o.ID,
			"first_name":  o.FirstName,
			"last_name":   o.LastName,
			"middle_name": o.MiddleName,
			"gender":      o.Gender,
			"phone":       o.Phone,
			"email":       o.Email,
			"password":    o.Password,
			"created_at":  o.CreatedAt,
			"updated_at":  o.UpdatedAt,
		}
	)

	_, err := u.session(ctx).NamedExecContext(ctx, insertUserQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert user: %w", err))
	}

	return nil
}

// UserByIDTx get user by id.
func (u *User) UserByIDTx(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var (
		getUserQuery = `
			SELECT 
				id, first_name, last_name, middle_name, gender, phone, email, created_at, updated_at, deleted_at
			FROM 
				users
			WHERE 
				id = ? AND
				deleted_at IS NULL`

		row UserRow
	)

	err := u.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getUserQuery), id)
	if err != nil {
		return domain.User{}, handleError(fmt.Errorf("failed to select user by id: %w", err))
	}

	return row.toDomain(), nil
}

// UserByEmailTx get user by email.
func (u *User) UserByEmailTx(ctx context.Context, email string) (domain.User, error) {
	var (
		getUserQuery = `
			SELECT 
				id, first_name, last_name, middle_name, password, gender, phone, email, created_at, updated_at, deleted_at
			FROM 
				users
			WHERE 
				email = ? AND
				deleted_at IS NULL`

		row UserRow
	)

	err := u.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getUserQuery), email)
	if err != nil {
		return domain.User{}, handleError(fmt.Errorf("failed to select user by id: %w", err))
	}

	return row.toDomainWithPassword(), nil
}

// UsersByIDsTx get users by ids.
func (u *User) UsersByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.Users, error) {
	var (
		getUserQuery = `
			SELECT 
				id, first_name, last_name, middle_name, gender, phone, email, created_at, updated_at, deleted_at
			FROM 
				users
			WHERE 
				id IN (?) AND
				deleted_at IS NULL`

		rows UserRows
	)

	if len(ids) == 0 {
		return domain.Users{}, nil
	}

	q, args, err := sqlx.In(getUserQuery, ids)
	if err != nil {
		return domain.Users{}, handleError(fmt.Errorf("failed to select users by ids: %w", err))
	}

	err = u.session(ctx).SelectContext(ctx, &rows, sqlx.Rebind(sqlx.DOLLAR, q), args...)
	if err != nil {
		return domain.Users{}, handleError(fmt.Errorf("failed to select users by ids: %w", err))
	}

	return rows.toDomain(), nil
}

// GetUserListTx gets a user list.
func (u *User) GetUserListTx(ctx context.Context, filters domain.UserListFilter) (domain.Users, error) {
	params, filtersQuery, anySlices := usersListFilter(filters)

	var (
		users = make(UserRows, 0)
		err   error
	)

	query := `
		SELECT DISTINCT ON (users.id)
			users.id, 
			users.first_name, 
			users.last_name, 
			users.middle_name, 
			users.gender, 
			users.phone, 
			users.email, 

			users.created_at, 
			users.updated_at,
			users.deleted_at
		FROM 
			users
		LEFT JOIN
			user_roles ON user_roles.user_id = users.id
		LEFT JOIN 
			schools ON schools.id = user_roles.school_id
		LEFT JOIN 	
			educational_organizations ON 
				educational_organizations.id = user_roles.organization_id OR 
				schools.organization_id = educational_organizations.id
	` + where(filtersQuery)

	query += fmt.Sprintf(
		` ORDER BY 
					users.id, users.created_at %s
				LIMIT ? OFFSET ?`,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	if anySlices {
		query, params, err = sqlx.In(query, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed on sqlx.In, err: %w", err))
		}
	}

	err = u.session(ctx).SelectContext(ctx, &users, sqlx.Rebind(sqlx.DOLLAR, query), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select user list: %w", err))
	}

	return users.toDomain(), nil
}

// UserCountTx gets a count of user.
func (u *User) UserCountTx(ctx context.Context, filters domain.UserListFilter) (int, error) {
	var (
		count int
		err   error
	)

	params, filtersQuery, anySlices := usersListFilter(filters)

	query := `
		SELECT
			COUNT(DISTINCT users.id)
		FROM 
			users
		LEFT JOIN
			user_roles ON user_roles.user_id = users.id
		LEFT JOIN 
			schools ON schools.id = user_roles.school_id
		LEFT JOIN 	
			educational_organizations ON 
				educational_organizations.id = user_roles.organization_id OR 
				schools.organization_id = educational_organizations.id
		` + where(filtersQuery)

	if anySlices {
		query, params, err = sqlx.In(query, params...)
		if err != nil {
			return 0, fmt.Errorf("failed on sqlx.In, err: %w", err)
		}
	}

	err = u.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, query), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get user list count: %w", err))
	}

	return count, nil
}

// usersListFilter returns query by user list filter.
func usersListFilter(filters domain.UserListFilter) (params []any, filtersQuery []string, anySlices bool) {
	anySlices = false

	if len(filters.OrganizationIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "educational_organizations.id IN (?)")
		params = append(params, filters.OrganizationIDs)
	}

	if len(filters.SchoolIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "schools.id IN (?)")
		params = append(params, filters.SchoolIDs)
	}

	if len(filters.Roles) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "user_roles.role IN (?)")
		params = append(params, filters.Roles)
	}

	if len(filters.Emails) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "users.email IN (?)")
		params = append(params, filters.Emails)
	}

	if len(filters.Phones) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "users.phone IN (?)")
		params = append(params, filters.Phones)
	}

	filtersQuery = append(filtersQuery, "users.deleted_at IS NULL")

	return params, filtersQuery, anySlices
}
