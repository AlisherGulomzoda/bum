package populate

import (
	"bum-service/internal/domain"
	"bum-service/internal/service/director"
	eduorganization "bum-service/internal/service/edu-organization"
	"bum-service/internal/service/headmaster"
	"bum-service/internal/service/owner"
	"bum-service/internal/service/school"
	"bum-service/internal/service/student"
	"bum-service/internal/service/teacher"
	"bum-service/internal/service/user"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
)

// runPopulate runs the Server.
func (s *Service) generateUserWithMultipleRoles(ctx context.Context) (domain.User, error) {
	// Add a user
	userDomain, err := s.container.Service.userService.AddUser(ctx,
		user.AddUserArgs{
			FirstName:  "FirstName",
			LastName:   "LastName",
			MiddleName: pointerFromString("MiddleName"),
			Gender:     string(domain.UserGenderTypeMale),
			Phone:      pointerFromString("+992900000000"),
			Email:      "bum@gmail.com",
			Password:   "test",
		},
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to add user: %w", err)
	}

	// Add a new organization
	organizationDomain, err := s.container.Service.eduOrganizationService.CreateEduOrganization(
		ctx, eduorganization.CreateEduOrganizationArgs{
			Name: gofakeit.Company(),
			Logo: pointerFromString(gofakeit.URL()),
		})

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create edu organization: %w", err)
	}

	// Add owner.
	_, err = s.container.Service.ownerService.AddOwner(ctx, owner.AddOwnerArgs{
		UserID:         userDomain.ID,
		OrganizationID: organizationDomain.ID,
		Phone:          pointerFromString(gofakeit.Phone()),
		Email:          pointerFromString(gofakeit.Email()),
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to add owner: %w", err)
	}

	// Add schools to organization
	schoolDomain, err := s.container.Service.schoolService.AddSchool(ctx, school.CreateSchoolArgs{
		Name:           gofakeit.School(),
		OrganizationID: organizationDomain.ID,
		Location:       gofakeit.Address().Address,
		Phone:          pointerFromString(gofakeit.Phone()),
		Email:          pointerFromString(gofakeit.Email()),
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to add school: %w", err)
	}

	// Add director
	_, err = s.container.Service.directorService.AddDirector(ctx, director.AddDirectorArgs{
		UserID:   userDomain.ID,
		SchoolID: schoolDomain.ID,
		Phone:    pointerFromString(gofakeit.Phone()),
		Email:    pointerFromString(gofakeit.Email()),
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to add director: %w", err)
	}

	// Add headmaster
	_, err = s.container.Service.headmasterService.AddHeadmaster(ctx, headmaster.AddHeadmasterArgs{
		UserID:   userDomain.ID,
		SchoolID: schoolDomain.ID,
		Phone:    pointerFromString(gofakeit.Phone()),
		Email:    pointerFromString(gofakeit.Email()),
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to add headmaster: %w", err)
	}

	groupDomain, err := s.container.Service.groupService.CreateGroup(ctx, school.CreateGroupArgs{
		SchoolID: schoolDomain.ID,
		Name:     "A",
		GradeID:  s.gradesList[TajikStandard].Grades[FifthGrade].ID,
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to add group: %w", err)
	}

	// Add student
	_, err = s.container.Service.studentService.AddStudent(ctx, student.AddStudentArgs{
		UserID:   userDomain.ID,
		SchoolID: schoolDomain.ID,
		GroupID:  groupDomain.ID,
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to add student: %w", err)
	}

	// Add Teacher
	_, err = s.container.Service.teacherService.AddTeacher(ctx, teacher.AddTeacherArgs{
		UserID:   userDomain.ID,
		SchoolID: schoolDomain.ID,
		Phone:    pointerFromString(gofakeit.Phone()),
		Email:    pointerFromString(gofakeit.Email()),
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to add teacher: %w", err)
	}

	return userDomain, nil
}
