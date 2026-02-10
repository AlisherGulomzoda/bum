package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// AddStudyPlanArgs is a struct for adding study plan.
type AddStudyPlanArgs struct {
	ID          *uuid.UUID
	Title       string
	Description *string
	PlanOrder   int16
}

// AssignStudyPlans assigns school group study plan.
func (s Service) AssignStudyPlans(
	ctx context.Context,
	groupSubjectID uuid.UUID,
	args []AddStudyPlanArgs,
) (domain.StudyPlans, error) {
	newStudyPlans := make(domain.StudyPlans, 0, len(args))

	for _, arg := range args {
		newStudyPlans = append(newStudyPlans, domain.NewStudyPlan(
			arg.ID,
			groupSubjectID,
			arg.Title,
			arg.Description,
			arg.PlanOrder,
			s.now,
		))
	}

	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.StudyPlans{}, fmt.Errorf("failed to begin transaction : %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf(
				"failed to end transaction on add study plans: %w: %w", domain.ErrInternalServerError, errEnd,
			)
		}
	}(tx)

	err = s.schoolRepo.AssignStudyPlansTx(txCtx, groupSubjectID, newStudyPlans)
	if err != nil {
		return domain.StudyPlans{}, fmt.Errorf("failed to add study plans repo: %w", err)
	}

	storedStudyPlans, err := s.schoolRepo.StudyPlanListTx(txCtx, groupSubjectID)
	if err != nil {
		return domain.StudyPlans{}, fmt.Errorf("failed to get study plans: %w", err)
	}

	return storedStudyPlans, nil
}

// StudyPlanList returns a list of study plan.
func (s Service) StudyPlanList(
	ctx context.Context,
	groupSubjectID uuid.UUID,
) (domain.StudyPlans, error) {
	list, err := s.schoolRepo.StudyPlanListTx(ctx, groupSubjectID)
	if err != nil {
		return domain.StudyPlans{}, fmt.Errorf("failed to get study plans: %w", err)
	}

	return list, nil
}

// StudyPlanChangeStatus set study plan status.
func (s Service) StudyPlanChangeStatus(
	ctx context.Context,
	groupSubjectID,
	studyPlanID uuid.UUID,
	status string,
) error {
	err := domain.StudyPlanStatus(status).Validate()
	if err != nil {
		return fmt.Errorf("failed to validate study plan status: %w", err)
	}

	err = s.schoolRepo.StudyPlanChangeStatusTx(ctx, groupSubjectID, studyPlanID, status)
	if err != nil {
		return fmt.Errorf("failed to get study plan: %w", err)
	}

	return nil
}
