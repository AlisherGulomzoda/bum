package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// AddRoleToUser creates adds a new role to the user.
func (s Service) AddRoleToUser(
	ctx context.Context,
	userID uuid.UUID,
	role domain.Role,
	schoolID *uuid.UUID,
	organizationID *uuid.UUID,
) (newUserRole domain.UserRole, err error) {
	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.UserRole{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf("failed to end transaction: %w: %w", domain.ErrInternalServerError, errEnd)
		}
	}(tx)

	newUserRole, err = domain.NewUserRole(userID, role, schoolID, organizationID, s.nowFunc)
	if err != nil {
		return domain.UserRole{}, fmt.Errorf("failed to create new user role domain: %w", err)
	}

	err = s.userRepo.AddUserRoleTx(txCtx, newUserRole, role.IsWithOnConflict())
	if err != nil {
		return domain.UserRole{}, fmt.Errorf("failed add user role to database: %w", err)
	}

	return newUserRole, nil
}
