package populate

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/internal/service/school"
)

var CommonTajikGroupNames = []string{
	"А", "Б", "В", "Г", "Ғ", "Д", "Е", "Ё", "Ж", "З", "И", "Ӣ", "Й", "К", "Қ", "Л", "М", "Н", "О", "П", "Р", "С", "Т", "У", "Ў", "Ф", "Х", "Ҳ", "Ч", "Ҷ", "Ш", "Ъ", "Э", "Ю", "Я",
}

func (s *Service) generateGroups(
	ctx context.Context, organization domain.EduOrganization, school domain.School, count int,
) (domain.Groups, error) {
	var (
		groups domain.Groups
		err    error
	)

	if school.GradeStandardID == nil {
		return domain.Groups{}, fmt.Errorf("school grade standard ID is required")
	}

	switch s.getGradeStandardByID(*school.GradeStandardID).Name {
	case TajikStandard:
		groups, err = s.generateTajikGroups(ctx, school.ID, count)
		if err != nil {
			return domain.Groups{}, err
		}
	case UKStandard:
		groups, err = s.generateUKGroups(ctx, school.ID, count)
		if err != nil {
			return domain.Groups{}, err
		}
	default:
		return domain.Groups{}, fmt.Errorf("standard is not ready")
	}

	return groups, nil
}

//nolint:funlen // it's ok
func (s *Service) generateTajikGroups(ctx context.Context, schoolID uuid.UUID, count int) (domain.Groups, error) {
	var (
		groups     domain.Groups
		groupsList []school.CreateGroupArgs
	)

	maxEduYear := 11
	if count < 11 {
		maxEduYear = count
	}

	for gradeYear := 1; gradeYear <= maxEduYear; gradeYear++ {
		groupsPerGrade := count / maxEduYear

		for commonGroupNameIndex := 0; commonGroupNameIndex < groupsPerGrade; commonGroupNameIndex++ {
			groupsList = append(groupsList, school.CreateGroupArgs{
				SchoolID: schoolID,
				Name:     CommonTajikGroupNames[commonGroupNameIndex],
				GradeID:  s.getGradeByName(TajikStandard, int8(gradeYear)).ID,
			})
		}
	}

	if len(groupsList) == 0 {
		fmt.Println("count count ==>", count)
		panic("count")
	}

	for _, g := range groupsList {
		group, err := s.generateGroup(ctx, g)
		if err != nil {
			return domain.Groups{}, err
		}

		groups = append(groups, group)
	}

	if len(groups) == 0 {
		fmt.Println("count count ==>", count)
		fmt.Println("groupsList count ==>", len(groupsList))
		panic("count")
	}

	return groups, nil
}

var CommonUKGroupNames = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

//nolint:funlen // it's ok
func (s *Service) generateUKGroups(ctx context.Context, schoolID uuid.UUID, count int) (domain.Groups, error) {
	var (
		groups     domain.Groups
		groupsList []school.CreateGroupArgs
	)

	maxEduYear := 11
	if count < 11 {
		maxEduYear = count
	}

	for gradeYear := 1; gradeYear <= maxEduYear; gradeYear++ {
		groupsPerGrade := count / maxEduYear

		for commonGroupNameIndex := 0; commonGroupNameIndex < groupsPerGrade; commonGroupNameIndex++ {
			groupsList = append(groupsList, school.CreateGroupArgs{
				SchoolID: schoolID,
				Name:     CommonUKGroupNames[commonGroupNameIndex],
				GradeID:  s.getGradeByName(UKStandard, int8(gradeYear)).ID,
			})
		}
	}

	for _, g := range groupsList {
		group, err := s.generateGroup(ctx, g)
		if err != nil {
			return domain.Groups{}, err
		}

		groups = append(groups, group)
	}

	if len(groups) == 0 {
		fmt.Println("count count ==>", count)
		fmt.Println("groupsList count ==>", len(groupsList))
		panic("count")
	}

	return groups, nil
}

//nolint:unparam // it's ok
func (s *Service) generateGroup(ctx context.Context, args school.CreateGroupArgs) (domain.Group, error) {
	group, err := s.schoolService().CreateGroup(ctx, args)
	if err != nil {
		return domain.Group{}, fmt.Errorf("failed to create group: %w", err)
	}

	return group, nil
}
