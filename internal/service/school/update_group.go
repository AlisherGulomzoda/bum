package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// UpdateGroupArgs is a list of arguments to update group.
type UpdateGroupArgs struct {
	ID      uuid.UUID
	Name    string
	GradeID uuid.UUID

	ClassTeacherID         *uuid.UUID
	ClassPresidentID       *uuid.UUID
	DeputyClassPresidentID *uuid.UUID
}

// UpdateGroup updates group.
func (s Service) UpdateGroup(ctx context.Context, args UpdateGroupArgs) (domain.Group, error) {
	group, err := s.groupRepo.GroupByIDTx(ctx, args.ID)
	if err != nil {
		return domain.Group{}, fmt.Errorf("failed to get group by id from database: %w", err)
	}

	group.Update(
		args.Name, args.GradeID, args.ClassTeacherID, args.ClassPresidentID, args.DeputyClassPresidentID, s.now,
	)

	err = s.groupRepo.UpdateGroupTx(ctx, group)
	if err != nil {
		return domain.Group{}, fmt.Errorf("failed to update group to database: %w", err)
	}

	group, err = s.GroupByID(ctx, group.ID)
	if err != nil {
		return domain.Group{}, fmt.Errorf("failed to get group by id: %w", err)
	}

	return group, nil
}
