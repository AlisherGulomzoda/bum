package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// UserRoles get user by id.
func (s Service) UserRoles(ctx context.Context, userID uuid.UUID) (domain.UserRoles, error) {
	// Get all user roles.
	userRoles, err := s.userRepo.UserRolesByIDTx(ctx, userID)
	if err != nil {
		return domain.UserRoles{}, fmt.Errorf("failed to get the user roles by id from database: %w", err)
	}

	// Collect ids of schools and organizations.
	schoolIDs := userRoles.SchoolIDs()
	organizationIDs := userRoles.OrganizationIDs()

	// Get all school info.
	schoolsShortInfos, err := s.schoolService.SchoolShortByIDs(ctx, schoolIDs)
	if err != nil {
		return domain.UserRoles{}, fmt.Errorf("failed get user school short info from service: %w", err)
	}

	// Add also organizations from school list to get info also for them.
	organizationIDs = append(organizationIDs, schoolsShortInfos.OrganizationIDs()...)

	// Get all organization info.
	organizationInfos, err := s.organizationService.EduOrganizationsShortInfoByIDs(ctx, organizationIDs)
	if err != nil {
		return domain.UserRoles{}, fmt.Errorf("failed get user school short info from service: %w", err)
	}

	// Set school and organization info to user roles.
	userRoles.SetSchoolShortInfosAndOrganization(schoolsShortInfos, organizationInfos)

	// if the user has a student role, then get info about its group.
	if userRoles.IsStudent() {
		studentsInfo, err := s.studentRepo.StudentsByUserIDTx(ctx, userID)
		if err != nil {
			return domain.UserRoles{}, fmt.Errorf("failed get student by user id from database: %w", err)
		}

		groups, err := s.groupService.GroupsByIDs(ctx, studentsInfo.GroupIDs())
		if err != nil {
			return domain.UserRoles{}, fmt.Errorf("failed get groups info by ids from group service: %w", err)
		}

		userRoles.SetGroups(studentsInfo, groups)
	}

	return userRoles, nil
}

// UserRolesByIDs get user by ids.
func (s Service) UserRolesByIDs(ctx context.Context, userIDs []uuid.UUID) (domain.UserRoles, error) {
	userRoles, err := s.userRepo.UserRolesByIDsTx(ctx, userIDs)
	if err != nil {
		return domain.UserRoles{}, fmt.Errorf("failed to get the users roles by ids from database: %w", err)
	}

	schoolIDs := userRoles.SchoolIDs()
	organizationIDs := userRoles.OrganizationIDs()

	schoolsShortInfos, err := s.schoolService.SchoolShortByIDs(ctx, schoolIDs)
	if err != nil {
		return domain.UserRoles{}, fmt.Errorf("failed get user school short info from service: %w", err)
	}

	organizationIDs = append(organizationIDs, schoolsShortInfos.OrganizationIDs()...)

	organizationInfos, err := s.organizationService.EduOrganizationsShortInfoByIDs(ctx, organizationIDs)
	if err != nil {
		return domain.UserRoles{}, fmt.Errorf("failed get user school short info from service: %w", err)
	}

	userRoles.SetSchoolShortInfosAndOrganization(schoolsShortInfos, organizationInfos)

	// TODO: add student info if user has student role.

	return userRoles, nil
}
