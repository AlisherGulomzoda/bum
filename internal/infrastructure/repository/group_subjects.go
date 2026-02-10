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

// GroupSubjects is a group subjects repository.
type GroupSubjects struct {
	db      postgres.DB
	session func(ctx context.Context) postgres.DB
}

// NewGroupSubjects create a new group subjects repository instance.
func NewGroupSubjects(db postgres.DB, session transaction.SessionDB) *GroupSubjects {
	return &GroupSubjects{
		db:      db,
		session: session.DB,
	}
}

const (
	// GroupSubjectsUniqueKey is unique key for group subjects groupID and school subjectID.
	GroupSubjectsUniqueKey = "group_subjects_key"
	// GroupSubjectsSchoolSubjectIDFkey is foreign key for school subjects.
	GroupSubjectsSchoolSubjectIDFkey = "group_subjects_school_subject_id_fkey"
	// GroupSubjectsGroupIDFkey is foreign key for groups.
	GroupSubjectsGroupIDFkey = "group_subjects_group_id_fkey"
	// GroupSubjectsTeacherIDFKey is foreign key for groups to teacher.
	GroupSubjectsTeacherIDFKey = "group_subjects_teacher_id_fkey"
)

const (
	// GroupNameKey is uniq key for group name.
	GroupNameKey = "groups_name_key"
)

// GroupSubjectRow is a row of group subject.
type GroupSubjectRow struct {
	ID              uuid.UUID  `db:"id"`
	GroupID         uuid.UUID  `db:"group_id"`
	SchoolSubjectID uuid.UUID  `db:"school_subject_id"`
	TeacherID       *uuid.UUID `db:"teacher_id"`
	Count           *int16     `db:"count"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts group subject row to domain.
func (s GroupSubjectRow) toDomain() domain.GroupSubject {
	return domain.GroupSubject{
		ID:              s.ID,
		SchoolSubjectID: s.SchoolSubjectID,
		GroupID:         s.GroupID,
		TeacherID:       s.TeacherID,
		Count:           s.Count,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		DeletedAt:       s.DeletedAt,
	}
}

// GroupSubjectRows is a collection GroupSubjectRow.
type GroupSubjectRows []GroupSubjectRow

// toDomain converts SchoolSubjectRows to domain.
func (s GroupSubjectRows) toDomain() []domain.GroupSubject {
	list := make([]domain.GroupSubject, 0, len(s))

	for _, row := range s {
		list = append(list, row.toDomain())
	}

	return list
}

// AddGroupSubjectsTx adds subjects to a group.
func (g *GroupSubjects) AddGroupSubjectsTx(ctx context.Context, groupSubject domain.GroupSubject) error {
	insertQuery := `
			INSERT INTO group_subjects
	    		( id, group_id, school_subject_id, teacher_id, count, created_at, updated_at) 
			VALUES 
	    		(:id,:group_id,:school_subject_id,:teacher_id,:count,:created_at,:updated_at)`

	_, err := g.session(ctx).NamedExecContext(ctx, insertQuery, map[string]any{
		"id":                groupSubject.ID,
		"group_id":          groupSubject.GroupID,
		"school_subject_id": groupSubject.SchoolSubjectID,
		"teacher_id":        groupSubject.TeacherID,
		"count":             groupSubject.Count,

		"created_at": groupSubject.CreatedAt,
		"updated_at": groupSubject.UpdatedAt,
	})
	if err != nil {
		return handleError(fmt.Errorf("failed to assign group subjects : %w", err))
	}

	return nil
}

// GroupSubjectByIDTx returns group subject by id.
func (g *GroupSubjects) GroupSubjectByIDTx(ctx context.Context, id uuid.UUID) (domain.GroupSubject, error) {
	query := `
		SELECT 
		    id, school_subject_id, group_id, teacher_id, count, created_at, updated_at
		FROM 
		    group_subjects
		WHERE 
		    id = ? AND 
		    deleted_at IS NULL;`

	var (
		groupSubjectList GroupSubjectRow
		err              error
	)

	err = g.session(ctx).GetContext(ctx, &groupSubjectList, sqlx.Rebind(sqlx.DOLLAR, query), id)
	if err != nil {
		return domain.GroupSubject{}, handleError(fmt.Errorf("failed to get group subject list: %w", err))
	}

	return groupSubjectList.toDomain(), nil
}

// GroupSubjectListTx returns a list of subjects assigned to a group.
func (g *GroupSubjects) GroupSubjectListTx(ctx context.Context, groupID uuid.UUID) (domain.GroupSubjects, error) {
	query := `
		SELECT 
		    id, school_subject_id, group_id, teacher_id, count, created_at, updated_at
		FROM 
		    group_subjects
		WHERE 
		    group_id = ? AND 
		    deleted_at IS NULL;`

	var (
		groupSubjectList = make(GroupSubjectRows, 0)
		err              error
	)

	err = g.session(ctx).SelectContext(ctx, &groupSubjectList, sqlx.Rebind(sqlx.DOLLAR, query), groupID)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select group subject list: %w", err))
	}

	return groupSubjectList.toDomain(), nil
}
