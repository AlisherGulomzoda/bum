package director

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// AddDirectorArgs is arguments for adding a new director.
type AddDirectorArgs struct {
	UserID   uuid.UUID
	SchoolID uuid.UUID
	Phone    *string
	Email    *string
}

// AddDirector adds a new director.
func (s Service) AddDirector(ctx context.Context, arg AddDirectorArgs) (newDirector domain.Director, err error) {
	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.Director{}, fmt.Errorf("failed to begin transaction : %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf(
				"failed to end transaction on create director : %w: %w", domain.ErrInternalServerError, errEnd,
			)
		}
	}(tx)

	userRole, err := s.userService.AddRoleToUser(txCtx, arg.UserID, domain.RoleDirector, &arg.SchoolID, nil)
	if err != nil {
		return domain.Director{}, fmt.Errorf("failed to add role to user: %w", err)
	}

	newDirector = domain.NewDirector(
		userRole.ID,
		arg.UserID,
		arg.SchoolID,
		arg.Phone,
		arg.Email,
		s.now,
	)

	err = s.directorRepo.AddDirectorTx(txCtx, newDirector)
	if err != nil {
		return domain.Director{}, fmt.Errorf("failed create director to database : %w", err)
	}

	newDirector, err = s.DirectorByID(txCtx, newDirector.ID)
	if err != nil {
		return domain.Director{}, fmt.Errorf("failed get director by id : %w", err)
	}

	return newDirector, nil
}
