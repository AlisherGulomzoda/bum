package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// CreateGroupArgs is a list of arguments to create a new group.
type CreateGroupArgs struct {
	SchoolID uuid.UUID
	Name     string
	GradeID  uuid.UUID
}

// CreateGroup creates a new group.
func (s Service) CreateGroup(ctx context.Context, arg CreateGroupArgs) (domain.Group, error) {
	newGroupDomain := domain.NewGroup(arg.SchoolID, arg.Name, arg.GradeID, s.now)

	err := s.groupRepo.CreateGroupTx(ctx, newGroupDomain)
	if err != nil {
		return domain.Group{}, fmt.Errorf("failed to create a new group to database: %w", err)
	}

	responseGroupDomain, err := s.GroupByID(ctx, newGroupDomain.ID)
	if err != nil {
		return domain.Group{}, fmt.Errorf(
			"failed to get group by id group_id=%s: %w", newGroupDomain.ID.String(), err,
		)
	}

	return responseGroupDomain, nil
}
