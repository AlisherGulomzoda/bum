package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"bum-service/internal/domain"
)

const (
	// StudyPlanGroupSubjectPlanOrderUniqueKey is unique key for study plan groupSubjectID and plan order.
	StudyPlanGroupSubjectPlanOrderUniqueKey = "study_plan_group_subject_id_plan_order_key"
	// StudyPlanGroupSubjectIDFKey is foreign key for group subjects.
	StudyPlanGroupSubjectIDFKey = "study_plan_group_subject_id_fkey"
)

// AssignStudyPlansTx adds study plans.
func (s School) AssignStudyPlansTx(ctx context.Context, groupSubjectID uuid.UUID, studyPlans domain.StudyPlans) error {
	var (
		// TODO: добавить проверку на то есть ли занятия по данным планам и добавить select for update
		softDelete = `
			UPDATE 
			    study_plans
			SET 
			    deleted_at = now(), updated_at = now()
			WHERE 
			    group_subject_id = :group_subject_id;`

		insertQuery = `
			INSERT INTO study_plans
	    		( id, group_subject_id, title, description, plan_order, status, created_at, updated_at) 
			VALUES 
	    		(:id,:group_subject_id,:title,:description,:plan_order,:status,:created_at,:updated_at)
			ON CONFLICT 
				(id) 
			DO UPDATE SET 
				deleted_at = NULL, 
				updated_at = now(), 
				title = :title,
				description = :description, 
				plan_order = :plan_order, 
				status = :status;`
	)

	_, err := s.session(ctx).NamedExecContext(ctx, softDelete, map[string]any{"group_subject_id": groupSubjectID})
	if err != nil {
		return handleError(fmt.Errorf("failed to remove study plan : %w", err))
	}

	for _, studyPlan := range studyPlans {
		_, err = s.session(ctx).NamedExecContext(ctx, insertQuery, map[string]any{
			"id":               studyPlan.ID,
			"group_subject_id": studyPlan.GroupSubjectID,
			"title":            studyPlan.Title,
			"description":      studyPlan.Description,
			"plan_order":       studyPlan.PlanOrder,
			"status":           studyPlan.Status,
			"created_at":       studyPlan.CreatedAt,
			"updated_at":       studyPlan.UpdatedAt,
		})
		if err != nil {
			return handleError(fmt.Errorf("failed to add study plan : %w", err))
		}
	}

	return nil
}

// StudyPlanRow is a row of study plan.
type StudyPlanRow struct {
	ID             uuid.UUID `db:"id"`
	GroupSubjectID uuid.UUID `db:"group_subject_id"`
	Title          string    `db:"title"`
	Description    *string   `db:"description"`
	PlanOrder      int16     `db:"plan_order"`
	Status         string    `db:"status"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// toDomain converts group subject row to domain.
func (s StudyPlanRow) toDomain() domain.StudyPlan {
	return domain.StudyPlan{
		ID:             s.ID,
		GroupSubjectID: s.GroupSubjectID,
		Title:          s.Title,
		Description:    s.Description,
		PlanOrder:      s.PlanOrder,
		Status:         domain.StudyPlanStatus(s.Status),
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
		DeletedAt:      s.DeletedAt,
	}
}

// StudyPlanRows is a collection StudyPlanRow.
type StudyPlanRows []StudyPlanRow

// toDomain converts StudyPlanRows to domain.
func (s StudyPlanRows) toDomain() domain.StudyPlans {
	list := make(domain.StudyPlans, 0, len(s))

	for _, row := range s {
		list = append(list, row.toDomain())
	}

	return list
}

// StudyPlanListTx returns a list of study plan.
func (s School) StudyPlanListTx(ctx context.Context, groupSubjectID uuid.UUID) (domain.StudyPlans, error) {
	query := `
		SELECT 
		    id, group_subject_id, title, description, plan_order, status, created_at, updated_at
		FROM 
		    study_plans
		WHERE 
		    group_subject_id = ? AND deleted_at IS NULL
		ORDER BY 
		    plan_order;`

	var (
		studyPlanList = make(StudyPlanRows, 0)
		err           error
	)

	err = s.session(ctx).SelectContext(ctx, &studyPlanList, sqlx.Rebind(sqlx.DOLLAR, query), groupSubjectID)
	if err != nil {
		return nil, handleError(fmt.Errorf("failed to select study plan list: %w", err))
	}

	return studyPlanList.toDomain(), nil
}

// StudyPlanChangeStatusTx changes study plan status.
func (s School) StudyPlanChangeStatusTx(
	ctx context.Context,
	groupSubjectID,
	studyPlanID uuid.UUID,
	status string,
) error {
	query := `
		UPDATE 
		    study_plans
		SET 
		    status = :status, updated_at = now()
		WHERE 
		    id = :id AND group_subject_id = :group_subject_id AND deleted_at IS NULL`

	_, err := s.session(ctx).NamedExecContext(ctx, query, map[string]any{
		"id":               studyPlanID,
		"group_subject_id": groupSubjectID,
		"status":           status,
	})
	if err != nil {
		return handleError(fmt.Errorf("failed to update study plan status: %w", err))
	}

	return nil
}
