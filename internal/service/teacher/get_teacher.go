package teacher

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// TeacherByID get teacher by id.
func (s Service) TeacherByID(ctx context.Context, id uuid.UUID) (domain.Teacher, error) {
	teacher, err := s.teacherRepo.TeacherByIDTx(ctx, id)
	if err != nil {
		return teacher, fmt.Errorf("failed to get teacher by id: %w", err)
	}

	user, err := s.userInfoService.UserByID(ctx, teacher.UserID)
	if err != nil {
		return teacher, fmt.Errorf("failed to get user by id: %w", err)
	}

	teacher.SetUser(user)

	return teacher, nil
}

// TeachersByIDs get teachers by ids.
func (s Service) TeachersByIDs(ctx context.Context, ids []uuid.UUID) (domain.Teachers, error) {
	teachers, err := s.teacherRepo.TeachersByIDsTx(ctx, ids)
	if err != nil {
		return teachers, fmt.Errorf("failed to get teachers by ids: %w", err)
	}

	userIDs := teachers.UserIDs()

	users, err := s.userInfoService.UsersByIDs(ctx, userIDs)
	if err != nil {
		return teachers, fmt.Errorf("failed to get user by id: %w", err)
	}

	teachers.SetUsers(users)

	return teachers, nil
}
