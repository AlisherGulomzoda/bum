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

// Group is a group repository.
type Group struct {
	db      postgres.DB
	session func(ctx context.Context) postgres.DB
}

// NewGroup create a new group repository instance.
func NewGroup(db postgres.DB, session transaction.SessionDB) *Group {
	return &Group{
		db:      db,
		session: session.DB,
	}
}

const (
	// GroupsGradeIDFkey is foreign key for groups.
	GroupsGradeIDFkey = "groups_grade_id_fkey"
	// GroupsSchoolIDFkey is foreign key for groups.
	GroupsSchoolIDFkey = "groups_school_id_fkey"
	// GroupsClassTeacherIDFKey is foreign key for groups to teacher.
	GroupsClassTeacherIDFKey = "groups_class_teacher_id_fkey"
	// GroupsClassPresidentIDFKey is foreign key for groups to student.
	GroupsClassPresidentIDFKey = "groups_class_president_id_fkey"
)

// GroupRow represents a row of Group.
type GroupRow struct {
	ID       uuid.UUID `db:"id"`
	SchoolID uuid.UUID `db:"school_id"`
	Name     string    `db:"name"`
	GradeID  uuid.UUID `db:"grade_id"`

	ClassTeacherID         *uuid.UUID `db:"class_teacher_id"`
	ClassPresidentID       *uuid.UUID `db:"class_president_id"`
	DeputyClassPresidentID *uuid.UUID `db:"deputy_class_president_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts an object into a domain model.
func (e GroupRow) toDomain() domain.Group {
	return domain.Group{
		ID:       e.ID,
		SchoolID: e.SchoolID,
		Name:     e.Name,
		GradeID:  e.GradeID,

		ClassTeacherID:         e.ClassTeacherID,
		ClassPresidentID:       e.ClassPresidentID,
		DeputyClassPresidentID: e.DeputyClassPresidentID,

		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

// GroupRows is list of GroupRow.
type GroupRows []GroupRow

// toDomain converts an object into a domain model.
func (e GroupRows) toDomain() domain.Groups {
	list := make(domain.Groups, 0, len(e))

	for _, row := range e {
		list = append(list, row.toDomain())
	}

	return list
}

// CreateGroupTx creates a new group.
func (g Group) CreateGroupTx(ctx context.Context, o domain.Group) error {
	var (
		sqlQuery = `
			INSERT INTO groups
	    		( id, school_id, name, grade_id, created_at, updated_at) 
			VALUES 
	    		(:id,:school_id,:name,:grade_id,:created_at,:updated_at)`

		args = map[string]any{
			"id":        o.ID,
			"school_id": o.SchoolID,
			"name":      o.Name,
			"grade_id":  o.GradeID,

			"created_at": o.CreatedAt,
			"updated_at": o.UpdatedAt,
		}
	)

	_, err := g.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to insert group: %w", err))
	}

	return nil
}

// GroupByIDTx get group by id.
func (g Group) GroupByIDTx(ctx context.Context, id uuid.UUID) (domain.Group, error) {
	var (
		getHeadmasterQuery = `
			SELECT 
	    		id, school_id, name, grade_id, class_teacher_id, class_president_id, created_at, updated_at, deleted_at 			
			FROM 
				groups
			WHERE 
				id = ? AND 
				deleted_at IS NULL`

		row GroupRow
	)

	err := g.session(ctx).GetContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getHeadmasterQuery), id)
	if err != nil {
		return domain.Group{}, handleError(fmt.Errorf("failed to get group by id: %w", err))
	}

	return row.toDomain(), nil
}

// UpdateGroupTx updates group.
func (g Group) UpdateGroupTx(ctx context.Context, group domain.Group) error {
	var (
		sqlQuery = `
			UPDATE 
				groups
			SET
	    		id 							=:id, 
	    		school_id 					=:school_id, 
	    		name 						=:name, 
	    		grade_id 					=:grade_id, 
	    		class_teacher_id 			=:class_teacher_id,
	    		class_president_id 			=:class_president_id,
	    		deputy_class_president_id 	=:deputy_class_president_id,
	     
	    		updated_at 					=:updated_at 
			WHERE
			    id = :id AND
			    deleted_at IS NULL`

		args = map[string]any{
			"id":        group.ID,
			"school_id": group.SchoolID,
			"name":      group.Name,
			"grade_id":  group.GradeID,

			"class_teacher_id":          group.ClassTeacherID,
			"class_president_id":        group.ClassPresidentID,
			"deputy_class_president_id": group.DeputyClassPresidentID,

			"updated_at": group.UpdatedAt,
		}
	)

	_, err := g.session(ctx).NamedExecContext(ctx, sqlQuery, args)
	if err != nil {
		return handleError(fmt.Errorf("failed to update group by id: %w", err))
	}

	return nil
}

// GroupsByIDsTx get groups by ids.
func (g Group) GroupsByIDsTx(ctx context.Context, ids []uuid.UUID) (domain.Groups, error) {
	var (
		sqlQuery = `
			SELECT 
	    		id, 
	    		school_id, 
	    		name, 
	    		grade_id, 
	    		class_teacher_id, 
	    		class_president_id, 
	    		deputy_class_president_id, 
	    		
	    		created_at, 
	    		updated_at, 
	    		deleted_at 			
			FROM 
				groups
			WHERE 
				id IN (?) AND
				deleted_at IS NULL`

		rows GroupRows
	)

	if len(ids) == 0 {
		return nil, nil
	}

	sqlQuery, params, err := sqlx.In(sqlQuery, ids)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select groups by ids: %w", err))
	}

	err = g.session(ctx).SelectContext(ctx, &rows, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), params...)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to get group by id: %w", err))
	}

	return rows.toDomain(), nil
}

// GroupListTx get group list.
func (g Group) GroupListTx(
	ctx context.Context, schoolID uuid.UUID, _ domain.GroupFilters,
) (domain.Groups, error) {
	var (
		getHeadmasterQuery = `
			SELECT 
	    		id, 
	    		school_id, 
	    		name, 
	    		grade_id, 
	    		class_teacher_id, 
	    		class_president_id, 
	    		deputy_class_president_id, 
	    		
	    		created_at, 
	    		updated_at, 
	    		deleted_at 			
			FROM 
				groups
			WHERE 
				school_id = ? AND
				deleted_at IS NULL
			ORDER BY 
			    name`

		row GroupRows
	)

	err := g.session(ctx).SelectContext(ctx, &row, sqlx.Rebind(sqlx.DOLLAR, getHeadmasterQuery), schoolID)
	if err != nil {
		return domain.Groups{}, handleError(fmt.Errorf("failed to get group list by school_id: %w", err))
	}

	return row.toDomain(), nil
}

// GroupListCountTx get group list count.
func (g Group) GroupListCountTx(
	ctx context.Context, schoolID uuid.UUID, _ domain.GroupFilters,
) (int, error) {
	var (
		sqlQuery = `
			SELECT
				count(*)
			FROM 
				groups
			WHERE 
				school_id = ? AND
				deleted_at IS NULL
	`

		count int
	)

	err := g.session(ctx).QueryRowxContext(ctx, sqlx.Rebind(sqlx.DOLLAR, sqlQuery), schoolID).Scan(&count)
	if err != nil {
		return 0, handleError(fmt.Errorf("failed to get groups count by school_id: %w", err))
	}

	return count, nil
}
