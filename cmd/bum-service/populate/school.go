//nolint:gosec,gofumpt,wsl,gocognit //it's ok for generating test data.
package populate

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/internal/service/school"
)

func (s *Service) generateSchools(ctx context.Context, organization domain.EduOrganization) error {
	var (
		schoolCount = randomIntFromInterval(schoolPerOrganizationMinCount, schoolPerOrganizationMaxCount)
	)

	for _ = range schoolCount {
		// generate school
		createdSchool, err := s.generateSchoolAndFill(ctx, organization)
		if err != nil {
			return err
		}

		_ = createdSchool
	}

	return nil
}

//nolint:funlen,gocognit // it's ok cause the main logic is here.
func (s *Service) generateSchoolAndFill(
	ctx context.Context, organization domain.EduOrganization,
) (domain.School, error) {
	var (
		createdSchool domain.School
		err           error
	)

	createdSchool, err = s.generateSchool(ctx, organization.ID)
	if err != nil {
		return domain.School{}, err
	}

	// set grade standard.
	standard := TajikStandard

	if rand.Int()%10 != 0 {
		standard = UKStandard
	}

	standardID := s.gradesList[standard].ID

	createdSchool, err = s.schoolService().UpdateSchool(ctx, school.UpdateSchoolArgs{
		ID:              createdSchool.ID,
		Name:            createdSchool.Name,
		Location:        createdSchool.Location,
		Phone:           createdSchool.Phone,
		Email:           createdSchool.Email,
		GradeStandardID: &standardID,
	})
	if err != nil {
		return domain.School{}, err
	}

	err = s.fillSchool(ctx, organization, createdSchool)
	if err != nil {
		return domain.School{}, err
	}

	return createdSchool, nil
}

func (s *Service) fillSchool(
	ctx context.Context,
	organization domain.EduOrganization,
	createdSchool domain.School,
) error {
	// Adding Director.
	for try := 1; try <= 10; try++ {
		user := s.getRandomUser()
		userID := user.ID

		_, err := s.generateDirector(ctx, createdSchool.ID, userID)
		if err != nil {
			// if a user already has this role, then we need to change our user or create a new user.
			if errors.Is(err, domain.ErrUserRoleInSchoolAndOrganizationAlreadyExists) {
				continue
			}

			return err
		}

		break
	}

	// Adding Headmaster.
	for try := 1; try <= 10; try++ {
		user := s.getRandomUser()
		userID := user.ID

		_, err := s.generateHeadmaster(ctx, createdSchool.ID, userID)
		if err != nil {
			// if a user already has this role, then we need to change our user or create a new user.
			if errors.Is(err, domain.ErrUserRoleInSchoolAndOrganizationAlreadyExists) {
				continue
			}

			return err
		}

		break
	}

	// Adding Students.
	studentCount := randomIntFromInterval(studentsPerSchoolMinCount, studentsPerSchoolMaxCount)

	// Adding groups.
	studentsPerGroupsCount := randomIntFromInterval(studentsPerGroupMinCount, studentsPerGroupMaxCount)
	groupCount := studentCount / studentsPerGroupsCount

	studentCount = groupCount * studentsPerGroupsCount

	if groupCount == 0 {
		fmt.Println("studentsPerGroupsCount = ", studentsPerGroupsCount)
		fmt.Println("studentCount = ", studentCount)
		panic("wefwe")
	}

	groups, err := s.generateGroups(ctx, organization, createdSchool, groupCount)
	if err != nil {
		return err
	}

	_ = groups
	err = s.generateStudents(ctx, createdSchool.ID, groups, studentCount)
	if err != nil {
		return err
	}

	// Adding school subjects for the school.
	for _, globalSubject := range s.subjectsList {
		_, err := s.schoolService().CreateSchoolSubject(
			ctx, school.CreateSchoolSubjectArgs{
				SchoolID:    createdSchool.ID,
				SubjectID:   globalSubject.ID,
				Name:        globalSubject.Name,
				Description: globalSubject.Description,
			},
		)
		if err != nil {
			return err
		}
	}

	// Adding auditoriums.
	for i := 0; i < 50; i++ {
		_, err := s.schoolService().CreateAuditorium(ctx, createdSchool.ID, school.CreateAuditoriumArgs{
			Name: fmt.Sprintf("auditorium #%d", i),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:funlen,gocognit // it's ok cause the main logic is here.
func (s *Service) generateSchool(
	ctx context.Context, organizationID uuid.UUID,
) (domain.School, error) {
	var (
		createdSchool domain.School
		err           error
	)

	for try := 1; try <= 10; try++ {
		var (
			name     = gofakeit.School()
			location = gofakeit.Address().Address
			phone    = gofakeit.Phone()
			email    = gofakeit.Email()
		)

		// fist we create createSchool.
		createdSchool, err = s.schoolService().AddSchool(ctx, school.CreateSchoolArgs{
			Name:           name,
			OrganizationID: organizationID,
			Location:       location,
			Phone:          &phone,
			Email:          &email,
		})
		if err != nil {
			switch {
			case // continue unless we create a unique school
				errors.Is(err, domain.ErrSchoolAlreadyExists),
				errors.Is(err, domain.ErrSchoolPhoneAlreadyExists),
				errors.Is(err, domain.ErrSchoolEmailAlreadyExists):
				continue
			}

			return domain.School{}, fmt.Errorf("failed to create school: %w", err)
		}

		break
	}

	return createdSchool, nil
}
