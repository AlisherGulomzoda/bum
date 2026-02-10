package student

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// AddStudentArgs is arguments for adding a new student.
type AddStudentArgs struct {
	UserID   uuid.UUID
	SchoolID uuid.UUID
	GroupID  uuid.UUID
}

// AddStudent adds student.
func (s Service) AddStudent(ctx context.Context, args AddStudentArgs) (newStudent domain.Student, err error) {
	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed to begin transaction : %w", err)
	}

	// TODO: добавить проверку совпадает ли SchoolID с школой внутри которого GroupID находится.

	// TODO: взять school id из группы

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf(
				"failed to end transaction on create student : %w: %w", domain.ErrInternalServerError, errEnd,
			)
		}
	}(tx)

	userRole, err := s.userService.AddRoleToUser(txCtx, args.UserID, domain.RoleStudent, &args.SchoolID, nil)
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed create user with role : %w", err)
	}

	newStudent = domain.NewStudent(
		userRole.ID,
		args.UserID,
		args.GroupID,

		s.now,
	)

	err = s.studentRepo.AddStudentTx(txCtx, newStudent)
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed add student to database : %w", err)
	}

	newStudent, err = s.StudentByID(txCtx, newStudent.ID)
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed get student by id : %w", err)
	}

	return newStudent, nil
}
