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

// Subject is subject repository.
type Subject struct {
	db      postgres.DB
	session func(ctx context.Context) postgres.DB
}

// NewSubject creates a new subject repository.
func NewSubject(db postgres.DB, session transaction.SessionDB) *Subject {
	return &Subject{
		db:      db,
		session: session.DB,
	}
}

const (
	// SubjectsNameUniqueKey is subject name unique key.
	SubjectsNameUniqueKey = "subjects_name_key"
)

// SubjectRow is a row containing subject.
type SubjectRow struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// SubjectRows are collection of SubjectRow.
type SubjectRows []SubjectRow

func (s SubjectRows) toDomain() domain.Subjects {
	list := make(domain.Subjects, 0, len(s))

	for index := range s {
		list = append(list, s[index].toModel())
	}

	return list
}

// toModel converts an object into domain model.
func (e SubjectRow) toModel() domain.Subject {
	return domain.Subject{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   e.DeletedAt,
	}
}

// CreateSubjectTx creates a new subject.
func (r Subject) CreateSubjectTx(ctx context.Context, subject domain.Subject) error {
	var (
		insertQuery = `
			INSERT INTO subjects
				( id, name, description, created_at, updated_at)
			VALUES 
				(:id,:name,:description,:created_at,:updated_at);`

		args = map[string]any{
			"id":          subject.ID,
			"name":        subject.Name,
			"description": subject.Description,
			"created_at":  subject.CreatedAt,
			"updated_at":  subject.UpdatedAt,
		}
	)

	_, err := r.session(ctx).NamedExecContext(ctx, insertQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert subject: %w", err))
	}

	return nil
}

// GetSubjectByIDTx gets a subject by ID.
func (r Subject) GetSubjectByIDTx(ctx context.Context, id uuid.UUID) (domain.Subject, error) {
	var (
		query = `
				SELECT 
					id, name, description, created_at, updated_at, deleted_at
				FROM 
					subjects
				WHERE 
					id = ? AND
					deleted_at IS NULL`

		subjectRow SubjectRow
	)

	err := r.session(ctx).GetContext(ctx, &subjectRow, sqlx.Rebind(sqlx.DOLLAR, query), id)
	if err != nil {
		return domain.Subject{}, handleError(fmt.Errorf("failed to get subject by ID: %w", err))
	}

	return subjectRow.toModel(), nil
}

// GetSubjectListTx gets a subject list.
func (r Subject) GetSubjectListTx(ctx context.Context, filters domain.SubjectListFilter) (domain.Subjects, error) {
	var (
		subjects = make(SubjectRows, 0)

		query = fmt.Sprintf(`
			SELECT 
				id, name, description, created_at, updated_at, deleted_at
			FROM 
				subjects
			WHERE 
				deleted_at IS NULL
			ORDER BY 
				created_at %s
			LIMIT ? OFFSET ?
			`, filters.SortOrder)
	)

	err := r.session(ctx).SelectContext(ctx, &subjects, sqlx.Rebind(sqlx.DOLLAR, query), filters.Limit, filters.Offset)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to get subject list: %w", err))
	}

	return subjects.toDomain(), nil
}

// SubjectCountTx gets a count of subject.
func (r Subject) SubjectCountTx(ctx context.Context) (int, error) {
	var (
		query = `
			SELECT 
				count(*)
			FROM 
				subjects
			WHERE 
				deleted_at IS NULL`
		count int
	)

	err := r.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, query)).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get subject count: %w", err))
	}

	return count, nil
}
