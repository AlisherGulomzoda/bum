package populate

import (
	"bum-service/internal/domain"
	"bum-service/internal/service/student"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (s *Service) generateStudents(ctx context.Context, schoolID uuid.UUID, groups domain.Groups, count int) error {
	perGroup := count / len(groups)

	for _, group := range groups {

		for i := 0; i < perGroup; i++ {
			generateStudent, err := s.generateStudent(ctx, schoolID, group.ID, nil)
			if err != nil {
				return err
			}

			_ = generateStudent
		}
	}

	return nil
}

//nolint:unused //it will be used.
func (s *Service) generateStudent(
	ctx context.Context, schoolID uuid.UUID, groupID uuid.UUID, userID *uuid.UUID,
) (domain.Student, error) {
	var (
		createdStudent domain.Student
		err            error
	)

	if userID == nil {
		user, err := s.genUser(ctx)
		if err != nil {
			return domain.Student{}, err
		}

		userID = &user.ID
	}

	createdStudent, err = s.studentService().AddStudent(ctx, student.AddStudentArgs{
		UserID:   *userID,
		SchoolID: schoolID,
		GroupID:  groupID,
	})
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed to create student: %w", err)
	}

	return createdStudent, nil
}
