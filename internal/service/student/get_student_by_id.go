package student

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// StudentByID gets student by id.
func (s Service) StudentByID(ctx context.Context, studentID uuid.UUID) (domain.Student, error) {
	studentDomain, err := s.studentRepo.StudentByIDTx(ctx, studentID)
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed to get student by id from database: %w", err)
	}

	user, err := s.userInfoService.UserByID(ctx, studentDomain.UserID)
	if err != nil {
		err = fmt.Errorf("failed to get user by id: %w", err)
		return studentDomain, err
	}

	schoolShort, err := s.schoolService.SchoolShortByID(ctx, studentDomain.SchoolID)
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed get user school short info from service: %w", err)
	}

	group, err := s.groupService.GroupByID(ctx, studentDomain.GroupID)
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed get group info by id from group service: %w", err)
	}

	studentDomain.SetShortSchool(schoolShort)
	studentDomain.SetUser(user)
	studentDomain.SetGroup(group)

	return studentDomain, nil
}

// StudentShortInfoByID gets student short info by id.
func (s Service) StudentShortInfoByID(ctx context.Context, studentID uuid.UUID) (domain.Student, error) {
	studentDomain, err := s.studentRepo.StudentByIDTx(ctx, studentID)
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed to get student by id from database: %w", err)
	}

	user, err := s.userInfoService.UserByID(ctx, studentDomain.UserID)
	if err != nil {
		err = fmt.Errorf("failed to get user by id: %w", err)
		return studentDomain, err
	}

	studentDomain.SetUser(user)

	return studentDomain, nil
}
