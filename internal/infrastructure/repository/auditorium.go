package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"bum-service/internal/domain"
)

// AuditoriumRow is a row of school auditorium.
type AuditoriumRow struct {
	ID              uuid.UUID  `db:"id"`
	SchoolID        uuid.UUID  `db:"school_id"`
	Name            string     `db:"name"`
	SchoolSubjectID *uuid.UUID `db:"school_subject_id"`
	Description     *string    `db:"description"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

// toDomain converts school subject row to domain.
func (s AuditoriumRow) toDomain() domain.Auditorium {
	return domain.Auditorium{
		ID:              s.ID,
		SchoolID:        s.SchoolID,
		Name:            s.Name,
		SchoolSubjectID: s.SchoolSubjectID,
		Description:     s.Description,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		DeletedAt:       s.DeletedAt,
	}
}

// AuditoriumRows is a collection AuditoriumRow.
type AuditoriumRows []AuditoriumRow

// toDomain converts SchoolSubjectRows to domain.
func (s AuditoriumRows) toDomain() domain.Auditoriums {
	list := make(domain.Auditoriums, 0, len(s))

	for _, row := range s {
		list = append(list, row.toDomain())
	}

	return list
}

const (
	// AuditoriumsNameUniqueKey name is unique for auditoriums.
	AuditoriumsNameUniqueKey = "auditoriums_name_key"
	// AuditoriumsSchoolSubjectIDFkey is foreign key for school subjects.
	AuditoriumsSchoolSubjectIDFkey = "auditoriums_school_subject_id_fkey"
	// AuditoriumsSchoolIDFkey is foreign key for schools.
	AuditoriumsSchoolIDFkey = "auditoriums_school_id_fkey"
)

// CreateAuditoriumTx creates a new auditorium.
func (s School) CreateAuditoriumTx(ctx context.Context, o domain.Auditorium) error {
	var (
		sqlQuery = `
			INSERT INTO auditoriums (id, school_id, name, school_subject_id, description, created_at, updated_at) 
			VALUES 
				(:id,:school_id,:name,:school_subject_id,:description,:created_at,:updated_at);`

		args = map[string]any{
			"id":                o.ID,
			"school_id":         o.SchoolID,
			"name":              o.Name,
			"school_subject_id": o.SchoolSubjectID,
			"description":       o.Description,
			"created_at":        o.CreatedAt,
			"updated_at":        o.UpdatedAt,
		}
	)

	_, err := s.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert subject: %w", err))
	}

	return nil
}

// AuditoriumByIDAndSchoolIDTx gets school auditorium by id and school id.
func (s School) AuditoriumByIDAndSchoolIDTx(ctx context.Context, id, schoolID uuid.UUID) (domain.Auditorium, error) {
	var (
		sqlQuery = `
			SELECT 
				id, school_id, name, school_subject_id, description, created_at, updated_at
			FROM 
				auditoriums
			WHERE 
				id = ? AND
				school_id = ? AND
				deleted_at IS NULL
	`
		auditorium AuditoriumRow
	)

	err := s.session(ctx).GetContext(ctx, &auditorium, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), id, schoolID)
	if err != nil {
		return domain.Auditorium{}, handleError(fmt.Errorf("failed to get school auditorium by id: %w", err))
	}

	return auditorium.toDomain(), nil
}

// AuditoriumListTx get school auditorium list.
func (s School) AuditoriumListTx(
	ctx context.Context, filters domain.AuditoriumListFilters,
) (domain.Auditoriums, error) {
	params, filtersQuery, anySlices := auditoriumListFilter(filters)

	sqlQuery := `
			SELECT 
				id, school_id, name, school_subject_id, description, created_at, updated_at
			FROM	
				auditoriums
	` + where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		schoolSubjectList = make(AuditoriumRows, 0)
		err               error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select school auditoriums list: %w", err))
		}
	}

	err = s.session(ctx).SelectContext(ctx, &schoolSubjectList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select school auditoriums list: %w", err))
	}

	return schoolSubjectList.toDomain(), nil
}

func auditoriumListFilter(
	filters domain.AuditoriumListFilters,
) (params []any, filtersQuery []string, anySlices bool) {
	filtersQuery = append(filtersQuery, "deleted_at IS NULL")
	anySlices = false

	filtersQuery = append(filtersQuery, "school_id = ?")
	params = append(params, filters.SchoolID)

	return params, filtersQuery, anySlices
}

// AuditoriumListCountTx returns list count of school auditorium by filter from database.
func (s School) AuditoriumListCountTx(ctx context.Context, filters domain.AuditoriumListFilters) (int, error) {
	params, filtersQuery, anySlices := auditoriumListFilter(filters)

	q := `
		SELECT 
			COUNT(*) 
		FROM	
			auditoriums
	` + where(filtersQuery)

	var (
		count int
		err   error
	)

	if anySlices {
		q, params, err = sqlx.In(q, params...)
		if err != nil {
			return count, handleError(fmt.Errorf("failed to get school auditoriums list count : %w", err))
		}
	}

	err = s.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get school auditoriums list count : %w", err))
	}

	return count, nil
}
