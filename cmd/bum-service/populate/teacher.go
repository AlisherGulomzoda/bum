package populate

import (
	"context"
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/internal/service/teacher"
)

//nolint:unused //it will be used.
func (s *Service) generateTeacher(
	ctx context.Context, schoolID uuid.UUID, userID uuid.UUID,
) (domain.Teacher, error) {
	var (
		createdTeacher domain.Teacher
		err            error
	)

	for {
		var (
			phone = gofakeit.Phone()
			email = gofakeit.Email()
		)

		createdTeacher, err = s.teacherService().AddTeacher(ctx, teacher.AddTeacherArgs{
			UserID:   userID,
			SchoolID: schoolID,
			Phone:    &phone,
			Email:    &email,
		})
		if err != nil {
			switch {
			case // continue unless we create a unique director
				errors.Is(err, domain.ErrUserEmailAlreadyExists),
				errors.Is(err, domain.ErrUserPhoneAlreadyExists):
				continue
			}

			return domain.Teacher{}, fmt.Errorf("failed to create teacher: %w", err)
		}

		break
	}

	return createdTeacher, nil
}
