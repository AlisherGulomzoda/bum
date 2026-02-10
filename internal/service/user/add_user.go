package user

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
	"bum-service/pkg/utils"
)

// AddUserArgs is args for creating user.
type AddUserArgs struct {
	FirstName  string
	LastName   string
	MiddleName *string
	Gender     string
	Phone      *string
	Email      string
	Password   string
}

// AddUser creates a new user.
func (s Service) AddUser(ctx context.Context, args AddUserArgs) (newUser domain.User, err error) {
	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf("failed to end transaction: %w: %w", domain.ErrInternalServerError, errEnd)
		}
	}(tx)

	passwordHash, err := utils.HashPassword(args.Password, s.passwordCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser, err = domain.NewUser(
		args.FirstName,
		args.LastName,
		args.MiddleName,
		args.Gender,
		args.Phone,
		args.Email,
		passwordHash,
		s.nowFunc,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create new user domain: %w", err)
	}

	err = s.userRepo.AddUserTx(txCtx, newUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed create user to database: %w", err)
	}

	newUser, err = s.UserByID(txCtx, newUser.ID)
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}
