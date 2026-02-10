package headmaster

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// AddHeadmasterArgs is arguments for adding a new headmaster.
type AddHeadmasterArgs struct {
	UserID   uuid.UUID
	SchoolID uuid.UUID
	Phone    *string
	Email    *string
}

// AddHeadmaster adds a new headmaster.
func (s Service) AddHeadmaster(
	ctx context.Context,
	arg AddHeadmasterArgs,
) (newHeadmaster domain.Headmaster, err error) {
	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.Headmaster{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf(
				"failed to end transaction on create headmaster: %w: %w", domain.ErrInternalServerError, errEnd,
			)
		}
	}(tx)

	userRole, err := s.userService.AddRoleToUser(txCtx, arg.UserID, domain.RoleHeadmaster, &arg.SchoolID, nil)
	if err != nil {
		return domain.Headmaster{}, fmt.Errorf("failed add role to user: %w", err)
	}

	newHeadmaster = domain.NewHeadmaster(
		userRole.ID,
		arg.UserID,
		arg.SchoolID,
		arg.Phone,
		arg.Email,
		s.now,
	)

	err = s.headmasterRepo.AddHeadmasterTx(txCtx, newHeadmaster)
	if err != nil {
		return domain.Headmaster{}, fmt.Errorf("failed create headmaster to database: %w", err)
	}

	newHeadmaster, err = s.HeadmasterByID(txCtx, newHeadmaster.ID)
	if err != nil {
		return domain.Headmaster{}, fmt.Errorf("failed get headmaster by id: %w", err)
	}

	return newHeadmaster, nil
}
