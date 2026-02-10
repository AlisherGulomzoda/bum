package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// GroupByID gets group by id.
func (s Service) GroupByID(ctx context.Context, groupID uuid.UUID) (domain.Group, error) {
	groupDomain, err := s.groupRepo.GroupByIDTx(ctx, groupID)
	if err != nil {
		return domain.Group{}, fmt.Errorf("failed to create a new group to database: %w", err)
	}

	grade, err := s.gradeService.GradeByID(ctx, groupDomain.GradeID)
	if err != nil {
		return domain.Group{}, fmt.Errorf("failed to get grade by id: %w", err)
	}

	if groupDomain.HasClassTeacher() {
		classTeacher, err := s.teacherService.TeacherByID(ctx, *groupDomain.ClassTeacherID)
		if err != nil {
			return domain.Group{}, fmt.Errorf("failed to get class teacher by id: %w", err)
		}

		groupDomain.SetClassTeacher(classTeacher)
	}

	if groupDomain.HasClassPresident() {
		classPresident, err := s.studentService.StudentShortInfoByID(ctx, *groupDomain.ClassPresidentID)
		if err != nil {
			return domain.Group{}, fmt.Errorf("failed to get class president by id: %w", err)
		}

		groupDomain.SetClassPresident(classPresident)
	}

	groupDomain.SetGrade(grade)

	return groupDomain, nil
}

// GroupsByIDs gets groups by ids.
func (s Service) GroupsByIDs(ctx context.Context, ids []uuid.UUID) (domain.Groups, error) {
	groups, err := s.groupRepo.GroupsByIDsTx(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new group to database: %w", err)
	}

	gradeIDs := groups.GradeIDs()

	grades, err := s.gradeService.GradesByIDs(ctx, gradeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get grades by ids: %w", err)
	}

	groups.SetGrades(grades)

	return groups, nil
}
