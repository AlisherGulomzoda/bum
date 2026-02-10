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

// School is schools repository.
type School struct {
	db      postgres.DB
	session func(context.Context) postgres.DB
}

// NewSchool create school repository instance.
func NewSchool(db postgres.DB, session transaction.SessionDB) *School {
	return &School{
		db:      db,
		session: session.DB,
	}
}

const (
	// SchoolsNameKey is schools name unique key.
	SchoolsNameKey = "schools_name_key"
	// SchoolsOrganizationIDFKey is schools organization id foreign key.
	SchoolsOrganizationIDFKey = "schools_organization_id_fkey"
	// SchoolSubjectsSchoolIDFKey is school subject school_id foreign key.
	SchoolSubjectsSchoolIDFKey = "school_subjects_school_id_fkey"
	// SchoolSubjectsSubjectIDFKey is school subject subjects_id foreign key.
	SchoolSubjectsSubjectIDFKey = "school_subjects_subject_id_fkey"
	// SchoolGradeStandardIDFKey is grade standard foreign key.
	SchoolGradeStandardIDFKey = "schools_grade_standard_id_fkey"
)

// SchoolRow represents a row of School.
type SchoolRow struct {
	ID              uuid.UUID  `db:"id"`
	Name            string     `db:"name"`
	OrganizationID  uuid.UUID  `db:"organization_id"`
	Location        string     `db:"location"`
	Phone           *string    `db:"phone"`
	Email           *string    `db:"email"`
	GradeStandardID *uuid.UUID `db:"grade_standard_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts an object into domain model.
func (e SchoolRow) toDomain() domain.School {
	return domain.School{
		ID:              e.ID,
		Name:            e.Name,
		OrganizationID:  e.OrganizationID,
		Location:        e.Location,
		Phone:           e.Phone,
		Email:           e.Email,
		GradeStandardID: e.GradeStandardID,

		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

// toShortDomain converts an object into domain model.
func (e SchoolRow) toShortDomain() domain.SchoolShortInfo {
	return domain.SchoolShortInfo{
		ID:             e.ID,
		Name:           e.Name,
		OrganizationID: e.OrganizationID,
	}
}

// SchoolRows are collection of SchoolRow.
type SchoolRows []SchoolRow

func (s SchoolRows) toDomain() domain.Schools {
	list := make(domain.Schools, 0, len(s))

	for index := range s {
		list = append(list, s[index].toDomain())
	}

	return list
}

func (s SchoolRows) toShortDomain() domain.SchoolShortInfos {
	list := make(domain.SchoolShortInfos, 0, len(s))

	for index := range s {
		list = append(list, s[index].toShortDomain())
	}

	return list
}

// CreateSchoolTx creates a new school.
func (s School) CreateSchoolTx(ctx context.Context, o domain.School) error {
	var (
		sqlQuery = `
			INSERT INTO schools
				( id, name, organization_id, location, phone, email, created_at, updated_at ) 
			VALUES 
				(:id,:name,:organization_id,:location,:phone,:email,:created_at,:updated_at ) 
	`
		args = map[string]any{
			"id":              o.ID,
			"name":            o.Name,
			"organization_id": o.OrganizationID,
			"location":        o.Location,
			"phone":           o.Phone,
			"email":           o.Email,

			"created_at": o.CreatedAt,
			"updated_at": o.UpdatedAt,
		}
	)

	_, err := s.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert school: %w", err))
	}

	return nil
}

// UpdateSchoolTx creates a new school.
func (s School) UpdateSchoolTx(ctx context.Context, o domain.School) error {
	var (
		sqlQuery = `
			UPDATE 
				schools
			SET	
			    name =:name, 
			    location =:location, 
			    phone =:phone, 
			    email =:email, 
			    grade_standard_id =:grade_standard_id,
			    
			    created_at =:created_at, 
			    updated_at =:updated_at 
			WHERE 
			    id =:id AND 
			    deleted_at IS NULL
	`
		args = map[string]any{
			"id":                o.ID,
			"name":              o.Name,
			"location":          o.Location,
			"phone":             o.Phone,
			"email":             o.Email,
			"grade_standard_id": o.GradeStandardID,

			"created_at": o.CreatedAt,
			"updated_at": o.UpdatedAt,
		}
	)

	_, err := s.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to update school: %w", err))
	}

	return nil
}

// SchoolByIDTx get a school by id.
func (s School) SchoolByIDTx(ctx context.Context, schoolID uuid.UUID) (domain.School, error) {
	var (
		sqlQuery = `
			SELECT 
				id, 
				name, 
				organization_id, 
				location, 
				phone, 
				email, 
				grade_standard_id, 
				
				created_at, 
				updated_at, 
				deleted_at
			FROM
				schools
			WHERE
				id = ? AND
				deleted_at IS NULL`

		school SchoolRow
	)

	if err := s.session(ctx).GetContext(ctx, &school, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), schoolID); err != nil {
		return domain.School{}, handleError(fmt.Errorf("failed to get school by ID: %w", err))
	}

	return school.toDomain(), nil
}

// SchoolListTx get school list.
func (s School) SchoolListTx(ctx context.Context, filters domain.SchoolFilters) (domain.Schools, error) {
	params, filtersQuery, anySlices := schoolsListFilter(filters)

	sqlQuery := `
			SELECT 
				id,
				name,
				organization_id,
				location,
				phone,
				email,
				grade_standard_id,

				created_at,
				updated_at,
				deleted_at
			FROM	
				schools
	` + where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		schoolList = make(SchoolRows, 0)
		err        error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select school list: %w", err))
		}
	}

	err = s.session(ctx).SelectContext(ctx, &schoolList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select school list: %w", err))
	}

	return schoolList.toDomain(), nil
}

func schoolsListFilter(filters domain.SchoolFilters) (params []any, filtersQuery []string, anySlices bool) {
	filtersQuery = append(filtersQuery, "deleted_at IS NULL")
	anySlices = false

	if len(filters.Emails) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "email IN (?)")
		params = append(params, filters.Emails)
	}

	if len(filters.Phones) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "phone IN (?)")
		params = append(params, filters.Phones)
	}

	if len(filters.OrganizationIDs) > 0 {
		anySlices = true

		filtersQuery = append(filtersQuery, "organization_id IN (?)")
		params = append(params, filters.OrganizationIDs)
	}

	return params, filtersQuery, anySlices
}

// SchoolListCountTx returns list count of school by filter from database.
func (s School) SchoolListCountTx(ctx context.Context, filters domain.SchoolFilters) (int, error) {
	params, filtersQuery, anySlices := schoolsListFilter(filters)

	q := `
		SELECT 
			COUNT(*) 
		FROM	
			schools
	` + where(filtersQuery)

	var (
		count int
		err   error
	)

	if anySlices {
		q, params, err = sqlx.In(q, params...)
		if err != nil {
			return count, handleError(fmt.Errorf("failed to get school list count : %w", err))
		}
	}

	err = s.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get school list count : %w", err))
	}

	return count, nil
}

// SchoolShortByIDTx get school short by id.
func (s School) SchoolShortByIDTx(ctx context.Context, id uuid.UUID) (domain.SchoolShortInfo, error) {
	var (
		sqlQuery = `
			SELECT 
				id, 
				name
			FROM	
				schools
			WHERE 
			    id = ? AND
			    deleted_at IS NULL`

		school SchoolRow
		err    error
	)

	err = s.session(ctx).GetContext(ctx, &school, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), id)
	if err != nil {
		return domain.SchoolShortInfo{}, handleError(fmt.Errorf("failed to get school by id: %w", err))
	}

	return school.toShortDomain(), nil
}

// SchoolShortByIDsTx get school short by ids.
func (s School) SchoolShortByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.SchoolShortInfos, error) {
	sqlQuery := `
			SELECT 
				id, 
				name,
				organization_id
			FROM	
				schools
			WHERE 
			    id IN(?)`

	var (
		schoolList = make(SchoolRows, 0)
		err        error
	)

	if len(ids) == 0 {
		return domain.SchoolShortInfos{}, nil
	}

	sqlQuery, params, err := sqlx.In(sqlQuery, ids)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select school list: %w", err))
	}

	err = s.session(ctx).SelectContext(ctx, &schoolList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select school list: %w", err))
	}

	return schoolList.toShortDomain(), nil
}

// SchoolSubjectRow is a row of school subject.
type SchoolSubjectRow struct {
	ID          uuid.UUID `db:"id"`
	SchoolID    uuid.UUID `db:"school_id"`
	SubjectID   uuid.UUID `db:"subject_id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts school subject row to domain.
func (s SchoolSubjectRow) toDomain() domain.SchoolSubject {
	return domain.SchoolSubject{
		ID:          s.ID,
		SubjectID:   s.SubjectID,
		SchoolID:    s.SchoolID,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		DeletedAt:   s.DeletedAt,
	}
}

// SchoolSubjectRows is a collection SchoolSubjectRow.
type SchoolSubjectRows []SchoolSubjectRow

// toDomain converts SchoolSubjectRows to domain.
func (s SchoolSubjectRows) toDomain() domain.SchoolSubjects {
	list := make(domain.SchoolSubjects, 0, len(s))

	for _, schoolSubjectRow := range s {
		list = append(list, schoolSubjectRow.toDomain())
	}

	return list
}

// CreateSchoolSubjectTx creates a new school subject.
func (s School) CreateSchoolSubjectTx(ctx context.Context, o domain.SchoolSubject) error {
	var (
		sqlQuery = `
			INSERT INTO school_subjects
				( id, subject_id, school_id, name, description, created_at, updated_at ) 
			VALUES 
				(:id,:subject_id,:school_id,:name,:description,:created_at,:updated_at ) 
	`
		args = map[string]any{
			"id":          o.ID,
			"subject_id":  o.SubjectID,
			"school_id":   o.SchoolID,
			"name":        o.Name,
			"description": o.Description,
			"created_at":  o.CreatedAt,
			"updated_at":  o.UpdatedAt,
		}
	)

	_, err := s.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert school subject: %w", err))
	}

	return nil
}

// SchoolSubjectByIDAndSchoolIDTx gets school subject by id and school id.
func (s School) SchoolSubjectByIDAndSchoolIDTx(
	ctx context.Context,
	id uuid.UUID,
	schoolID uuid.UUID,
) (domain.SchoolSubject, error) {
	var (
		sqlQuery = `
			SELECT 
				id, subject_id, school_id, name, description, created_at, updated_at
			FROM 
				school_subjects
			WHERE 
				id = ? AND
				school_id = ? AND
				deleted_at IS NULL
	`
		schoolSubject SchoolSubjectRow
	)

	err := s.session(ctx).GetContext(ctx, &schoolSubject, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), id, schoolID)
	if err != nil {
		return domain.SchoolSubject{}, handleError(fmt.Errorf("failed to get school subject by id: %w", err))
	}

	return schoolSubject.toDomain(), nil
}

// SchoolSubjectsByIDs get school subjects by ids.
func (s School) SchoolSubjectsByIDs(ctx context.Context, ids []uuid.UUID) (domain.SchoolSubjects, error) {
	var (
		sqlQuery = `
			SELECT 
				id, subject_id, school_id, name, description, created_at, updated_at
			FROM	
				school_subjects
			WHERE 
			    id IN (?) AND 
				deleted_at IS NULL	
	`

		schoolSubjectList = make(SchoolSubjectRows, 0)
		err               error
	)

	if len(ids) == 0 {
		return nil, nil
	}

	sqlQuery, params, err := sqlx.In(sqlQuery, ids)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select school subjects by ids: %w", err))
	}

	err = s.session(ctx).SelectContext(ctx, &schoolSubjectList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select school subjects by ids: %w", err))
	}

	return schoolSubjectList.toDomain(), nil
}

// SchoolSubjectListTx get school subject list.
func (s School) SchoolSubjectListTx(
	ctx context.Context, filters domain.SchoolSubjectFilters,
) (domain.SchoolSubjects, error) {
	params, filtersQuery, anySlices := schoolsSubjectListFilter(filters)

	sqlQuery := `
			SELECT 
				id, subject_id, school_id, name, description, created_at, updated_at
			FROM	
				school_subjects
	` + where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	var (
		schoolSubjectList = make(SchoolSubjectRows, 0)
		err               error
	)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select school subject list: %w", err))
		}
	}

	err = s.session(ctx).SelectContext(ctx, &schoolSubjectList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select school subject list: %w", err))
	}

	return schoolSubjectList.toDomain(), nil
}

func schoolsSubjectListFilter(
	filters domain.SchoolSubjectFilters,
) (params []any, filtersQuery []string, anySlices bool) {
	filtersQuery = append(filtersQuery, "deleted_at IS NULL")
	anySlices = false

	filtersQuery = append(filtersQuery, "school_id = ?")
	params = append(params, filters.SchoolID)

	return params, filtersQuery, anySlices
}

// SchoolSubjectsListCountTx returns list count of school subjects by filter from database.
func (s School) SchoolSubjectsListCountTx(ctx context.Context, filters domain.SchoolSubjectFilters) (int, error) {
	params, filtersQuery, anySlices := schoolsSubjectListFilter(filters)

	q := `
		SELECT 
			COUNT(*) 
		FROM	
			school_subjects
	` + where(filtersQuery)

	var (
		count int
		err   error
	)

	if anySlices {
		q, params, err = sqlx.In(q, params...)
		if err != nil {
			return count, handleError(fmt.Errorf("failed to get school subjects list count : %w", err))
		}
	}

	err = s.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get school subjects list count : %w", err))
	}

	return count, nil
}
