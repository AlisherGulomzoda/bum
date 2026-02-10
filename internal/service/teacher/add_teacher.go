package teacher

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// AddTeacherArgs is args for create teacher.
type AddTeacherArgs struct {
	UserID   uuid.UUID
	SchoolID uuid.UUID
	Phone    *string
	Email    *string
}

// AddTeacher adds a new teacher.
func (s Service) AddTeacher(ctx context.Context, args AddTeacherArgs) (newTeacher domain.Teacher, err error) {
	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.Teacher{}, fmt.Errorf("failed to begin transaction : %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf(
				"failed to end transaction on create teacher : %w: %w", domain.ErrInternalServerError, errEnd,
			)
		}
	}(tx)

	userRole, err := s.userService.AddRoleToUser(txCtx, args.UserID, domain.RoleTeacher, &args.SchoolID, nil)
	if err != nil {
		return domain.Teacher{}, fmt.Errorf("failed add role to user: %w", err)
	}

	newTeacher = domain.NewTeacher(
		userRole.ID,
		args.UserID,
		args.SchoolID,
		args.Phone,
		args.Email,
		s.now,
	)

	err = s.teacherRepo.CreateTeacherTx(txCtx, newTeacher)
	if err != nil {
		return domain.Teacher{}, fmt.Errorf("failed create teacher to database : %w", err)
	}

	newTeacher, err = s.TeacherByID(txCtx, newTeacher.ID)
	if err != nil {
		return domain.Teacher{}, fmt.Errorf("failed get teacher by id : %w", err)
	}

	return newTeacher, nil
}
