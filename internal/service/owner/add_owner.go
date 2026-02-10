package owner

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// AddOwnerArgs is arguments for adding a new owner.
type AddOwnerArgs struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Phone          *string
	Email          *string
}

// AddOwner adds a new owner.
func (s Service) AddOwner(ctx context.Context, arg AddOwnerArgs) (newOwner domain.Owner, err error) {
	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.Owner{}, fmt.Errorf("failed to begin transaction : %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf(
				"failed to end transaction on add owner : %w: %w", domain.ErrInternalServerError, errEnd,
			)
		}
	}(tx)

	userRole, err := s.userService.AddRoleToUser(txCtx, arg.UserID, domain.RoleOwner, nil, &arg.OrganizationID)
	if err != nil {
		return domain.Owner{}, fmt.Errorf("failed to add role to user: %w", err)
	}

	newOwner = domain.NewOwner(
		userRole.ID,
		arg.UserID,
		arg.OrganizationID,
		arg.Phone,
		arg.Email,
		s.now,
	)

	err = s.ownerRepo.AddOwnerTx(txCtx, newOwner)
	if err != nil {
		return domain.Owner{}, fmt.Errorf("failed add owner to database : %w", err)
	}

	newOwner, err = s.OwnerByID(txCtx, newOwner.ID)
	if err != nil {
		return domain.Owner{}, fmt.Errorf("failed get owner by id : %w", err)
	}

	return newOwner, nil
}
