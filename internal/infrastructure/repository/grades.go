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

// Grades is grades repository.
type Grades struct {
	db      postgres.DB
	session func(ctx context.Context) postgres.DB
}

// NewGrades create a new grades repository instance.
func NewGrades(db postgres.DB, session transaction.SessionDB) *Grades {
	return &Grades{
		db:      db,
		session: session.DB,
	}
}

const (
	// GradeStandardsNameUniqueKey is grade standards name unique key.
	GradeStandardsNameUniqueKey = "grade_standards_name_key"
	// GradesNameUniqueKey is grades name unique key.
	GradesNameUniqueKey = "grades_name_key"
	// GradeStandardsOrganizationIDFKey is grade standards organization_id foreign key.
	GradeStandardsOrganizationIDFKey = "grade_standards_organization_id_fkey"
	// GradesGradeStandardIDFKey is grades grade_standard_id foreign key.
	GradesGradeStandardIDFKey = "grades_grade_standard_id_fkey"
)

// GradeStandardRows is list of GradeStandardRow.
type GradeStandardRows []GradeStandardRow

// toDomain converts to entity.
func (g GradeStandardRows) toDomain() domain.GradeStandards {
	list := make(domain.GradeStandards, 0, len(g))

	for index := range g {
		list = append(list, g[index].toDomain())
	}

	return list
}

// GradeStandardRow is a row containing grade standard.
type GradeStandardRow struct {
	ID             uuid.UUID  `db:"id"`
	OrganizationID *uuid.UUID `db:"organization_id"`
	Name           string     `db:"name"`
	EducationYears int8       `db:"education_years"`
	Description    *string    `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts an object into domain model.
func (g GradeStandardRow) toDomain() domain.GradeStandard {
	return domain.GradeStandard{
		ID:             g.ID,
		OrganizationID: g.OrganizationID,
		Name:           g.Name,
		EducationYears: g.EducationYears,
		Description:    g.Description,
		CreatedAt:      g.CreatedAt,
		UpdatedAt:      g.UpdatedAt,
		DeletedAt:      g.DeletedAt,
	}
}

// GradeRows is a collection rows containing grades.
type GradeRows []GradeRow

// toDomain converts an object into domain model.
func (g GradeRows) toDomain() domain.Grades {
	list := make(domain.Grades, 0, len(g))

	for index := range g {
		list = append(list, g[index].toDomain())
	}

	return list
}

// GradeRow is a row containing grade.
type GradeRow struct {
	ID              uuid.UUID `db:"id"`
	GradeStandardID uuid.UUID `db:"grade_standard_id"`
	Name            string    `db:"name"`
	EducationYear   *int8     `db:"education_year"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts an object into domain model.
func (g GradeRow) toDomain() domain.Grade {
	return domain.Grade{
		ID:              g.ID,
		GradeStandardID: g.GradeStandardID,
		Name:            g.Name,
		EducationYear:   g.EducationYear,
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
		DeletedAt:       g.DeletedAt,
	}
}

// prepareGradeRows prepare grade rows from domain.
func prepareGradeRows(grades domain.Grades) GradeRows {
	rows := make(GradeRows, len(grades))

	for i := range grades {
		rows[i] = GradeRow{
			ID:              grades[i].ID,
			GradeStandardID: grades[i].GradeStandardID,
			Name:            grades[i].Name,
			EducationYear:   grades[i].EducationYear,
			CreatedAt:       grades[i].CreatedAt,
			UpdatedAt:       grades[i].UpdatedAt,
			DeletedAt:       grades[i].DeletedAt,
		}
	}

	return rows
}

// CreateGradesTx creates a new grades within a transaction session.
func (g *Grades) CreateGradesTx(ctx context.Context, grades domain.Grades) error {
	rows := prepareGradeRows(grades)

	sqlQuery := `
		INSERT INTO grades
			( id, grade_standard_id, name, education_year, created_at, updated_at)
		VALUES 
			(:id,:grade_standard_id,:name,:education_year,:created_at,:updated_at)`

	_, err := g.session(ctx).NamedExecContext(ctx, sqlQuery, rows)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert grades: %w", err))
	}

	return nil
}

// CreateGradeStandardTx creates a new grade standard within a transaction session.
func (g *Grades) CreateGradeStandardTx(ctx context.Context, gradeStandard domain.GradeStandard) error {
	row := GradeStandardRow{
		ID:             gradeStandard.ID,
		OrganizationID: gradeStandard.OrganizationID,
		Name:           gradeStandard.Name,
		EducationYears: gradeStandard.EducationYears,
		Description:    gradeStandard.Description,
		CreatedAt:      gradeStandard.CreatedAt,
		UpdatedAt:      gradeStandard.UpdatedAt,
		DeletedAt:      gradeStandard.DeletedAt,
	}

	sqlQuery := `
		INSERT INTO grade_standards 
			( id, organization_id, name, education_years, description, created_at, updated_at)
		VALUES 
			(:id,:organization_id,:name,:education_years,:description,:created_at,:updated_at)
	`

	_, err := g.session(ctx).NamedExecContext(ctx, sqlQuery, row)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert grade standard: %w", err))
	}

	return nil
}

// GradesByGradeStandardIDTx get grade standard grades.
func (g *Grades) GradesByGradeStandardIDTx(ctx context.Context, gradeStandardID uuid.UUID) (domain.Grades, error) {
	grades := make(GradeRows, 0)

	sqlQuery := `
		SELECT 
		    id, grade_standard_id, name, education_year, created_at, updated_at, deleted_at	
		FROM 
		    grades 
		WHERE 
		    deleted_at IS NULL AND 
		    grade_standard_id = ?;`

	err := g.session(ctx).SelectContext(ctx, &grades, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), gradeStandardID)
	if err != nil {
		return domain.Grades{}, handleError(fmt.Errorf("failed to select grades by grade standard id: %w", err))
	}

	return grades.toDomain(), nil
}

// GradeStandardByIDTx get a grade standard by id.
func (g *Grades) GradeStandardByIDTx(ctx context.Context, id uuid.UUID) (domain.GradeStandard, error) {
	var gradeStandard GradeStandardRow

	sqlQuery := `
	SELECT 
	    id, organization_id, name, education_years, description, created_at, updated_at, deleted_at
	FROM 
	    grade_standards 
	WHERE 
	    deleted_at IS NULL AND 
	    id = ?;
	`

	err := g.session(ctx).GetContext(ctx, &gradeStandard, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), id)
	if err != nil {
		return domain.GradeStandard{}, handleError(fmt.Errorf("failed to select grade standard by id: %w", err))
	}

	return gradeStandard.toDomain(), nil
}

// gradeStandardListFilter returns query by grade standard list filter.
func gradeStandardListFilter(
	_ domain.GradeStandardListFilter,
) (_ []any, filtersQuery []string, anySlices bool) {
	filtersQuery = append(filtersQuery, "deleted_at IS NULL")
	anySlices = false

	// TODO: add filters here.

	return []any{}, filtersQuery, anySlices
}

// GradeStandardListTx returns list of grade standard by filter from database.
func (g *Grades) GradeStandardListTx(
	ctx context.Context,
	filters domain.GradeStandardListFilter,
) (domain.GradeStandards, error) {
	var (
		gradeStandardList = make(GradeStandardRows, 0)
		err               error
	)

	params, filtersQuery, anySlices := gradeStandardListFilter(filters)

	sqlQuery := `
		SELECT
			id, organization_id, name, education_years, description, created_at, updated_at, deleted_at
		FROM 
			grade_standards 
		` +
		where(filtersQuery)

	sqlQuery += fmt.Sprintf(
		` ORDER BY created_at %s 
			LIMIT ? OFFSET ? `,
		filters.SortOrder,
	)

	params = append(params, filters.Limit, filters.Offset)

	if anySlices {
		sqlQuery, params, err = sqlx.In(sqlQuery, params...)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to select grade standard list: %w", err))
		}
	}

	err = g.session(ctx).SelectContext(ctx, &gradeStandardList, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select grade standard list: %w", err))
	}

	return gradeStandardList.toDomain(), nil
}

// GradeStandardListCountTx returns count of grade standard by filter from database.
func (g *Grades) GradeStandardListCountTx(ctx context.Context, filters domain.GradeStandardListFilter) (int, error) {
	params, filtersQuery, anySlices := gradeStandardListFilter(filters)

	q := `SELECT COUNT(*) FROM grade_standards ` + where(filtersQuery)

	var (
		count int
		err   error
	)

	if anySlices {
		q, params, err = sqlx.In(q, params...)
		if err != nil {
			return count, handleError(fmt.Errorf("failed to get grade standard list count : %w", err))
		}
	}

	err = g.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, q), params...).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get grade standard list count: %w", err))
	}

	return count, nil
}

// GradesBysIDTx get a grades by ids.
func (g *Grades) GradesBysIDTx(ctx context.Context, ids []uuid.UUID) (domain.Grades, error) {
	grades := make(GradeRows, 0)

	sqlQuery := `
		SELECT 
		    id, grade_standard_id, name, education_year, created_at, updated_at, deleted_at	
		FROM 
		    grades 
		WHERE
		    deleted_at IS NULL AND 
		    id IN (?)
		ORDER BY education_year`

	if len(ids) == 0 {
		return domain.Grades{}, nil
	}

	sqlQuery, params, err := sqlx.In(sqlQuery, ids)
	if err != nil {
		return domain.Grades{}, handleError(fmt.Errorf("failed to select grades by grade standard id: %w", err))
	}

	err = g.session(ctx).SelectContext(ctx, &grades, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return domain.Grades{}, handleError(fmt.Errorf("failed to select grades by grade standard id: %w", err))
	}

	return grades.toDomain(), nil
}

// GradeByIDTx get a grade by id.
func (g *Grades) GradeByIDTx(ctx context.Context, id uuid.UUID) (domain.Grade, error) {
	var (
		sqlQuery = `
		SELECT 
		    id, grade_standard_id, name, education_year, created_at, updated_at, deleted_at	
		FROM 
		    grades 
		WHERE
		    id = ? AND
		    deleted_at IS NULL`

		grade GradeRow
	)

	err := g.session(ctx).GetContext(ctx, &grade, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), id)
	if err != nil {
		return domain.Grade{}, handleError(fmt.Errorf("failed to get grade by grade id: %w", err))
	}

	return grade.toDomain(), nil
}
