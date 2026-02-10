package domain

import (
	"time"

	"github.com/google/uuid"

	"bum-service/pkg/utils"
)

// User is structure of user.
type User struct {
	ID         uuid.UUID
	FirstName  string
	LastName   string
	MiddleName *string
	Gender     Gender
	Phone      *string
	Email      string
	Password   string

	UserRoles UserRoles

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewUser create new User domain.
func NewUser(
	firsName string,
	lastName string,
	middleName *string,
	gender string,
	phone *string,
	email string,
	password string,
	nowFunc func() time.Time,
) (User, error) {
	userGender := Gender(gender)
	if ok := userGender.Validate(); !ok {
		return User{}, ErrUserGenderBadRequest
	}

	now := nowFunc()

	return User{
		ID:         uuid.New(),
		FirstName:  firsName,
		LastName:   lastName,
		MiddleName: utils.GetStrValueFromArgument(middleName),
		Gender:     userGender,
		Phone:      phone,
		Email:      email,
		Password:   password,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// SetRoles sets user roles.
func (u *User) SetRoles(roles UserRoles) {
	u.UserRoles = roles
}

// Users is a list of User.
type Users []User

// UserIDs returns list of user ids.
func (u Users) UserIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(u))

	for _, user := range u {
		list = append(list, user.ID)
	}

	return list
}

// SetRoles sets user roles to users.
func (u Users) SetRoles(userRoles UserRoles) {
	mapOfUserRoles := make(map[uuid.UUID]UserRoles, len(u))

	for _, userRole := range userRoles {
		mapOfUserRoles[userRole.UserID] = append(mapOfUserRoles[userRole.UserID], userRole)
	}

	for index := range u {
		u[index].SetRoles(mapOfUserRoles[u[index].ID])
	}
}

// Gender represents the user's gender.
type Gender string

const (
	// UserGenderTypeMale user gender type male.
	UserGenderTypeMale Gender = "male"
	// UserGenderTypeFemale user gender type female.
	UserGenderTypeFemale Gender = "female"
)

// Validate validate user gender.
func (g Gender) Validate() bool {
	switch g {
	case UserGenderTypeFemale, UserGenderTypeMale:
		return true
	}

	return false
}

// UserRole is structure of UserRole.
type UserRole struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Role   Role

	SchoolID        *uuid.UUID
	SchoolShortInfo *SchoolShortInfo

	OrganizationID           *uuid.UUID
	EduOrganizationShortInfo *EduOrganizationShortInfo

	GroupID *uuid.UUID
	Group   *Group

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// SetSchoolShortInfo sets school short info.
func (u *UserRole) SetSchoolShortInfo(school SchoolShortInfo) {
	u.SchoolShortInfo = &school
}

// SetOrganizationShortInfo sets organization short info.
func (u *UserRole) SetOrganizationShortInfo(organization EduOrganizationShortInfo) {
	u.EduOrganizationShortInfo = &organization
	u.OrganizationID = &organization.ID
}

// SetGroupShortInfo sets group info.
func (u *UserRole) SetGroupShortInfo(group Group) {
	u.Group = &group
	u.GroupID = &group.ID
}

// UserRoles is slice of UserRole.
type UserRoles []UserRole

// SchoolIDs returns the list of UserRoles schools ids.
func (u UserRoles) SchoolIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(u))

	for _, userRole := range u {
		if userRole.SchoolID != nil {
			list = append(list, *userRole.SchoolID)
		}
	}

	return list
}

// OrganizationIDs returns the list of UserRoles organization ids.
func (u UserRoles) OrganizationIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(u))

	for _, userRole := range u {
		if userRole.OrganizationID != nil {
			list = append(list, *userRole.OrganizationID)
		}
	}

	return list
}

func (u UserRoles) IsStudent() bool {
	for _, userRole := range u {
		if userRole.IsStudent() {
			return true
		}
	}

	return false
}

// SetGroups sets group info for students.
// TODO убрать костыль.
func (u UserRoles) SetGroups(studentsInfo Students, groups Groups) {
	studentsInfo.SetGroups(groups)

	mapOfStudentsBySchoolID := make(map[uuid.UUID]Student, len(u))
	for _, student := range studentsInfo {
		mapOfStudentsBySchoolID[student.SchoolID] = student
	}

	for index := range u {
		if u[index].IsStudent() && u[index].SchoolID != nil {
			u[index].SetGroupShortInfo(mapOfStudentsBySchoolID[*u[index].SchoolID].Group)
		}
	}
}

// SetSchoolShortInfosAndOrganization sets SchoolShortInfo and organization info if it exists.
func (u UserRoles) SetSchoolShortInfosAndOrganization(
	schools SchoolShortInfos,
	organizations EduOrganizationShortInfos,
) {
	mapOfSchools := make(map[uuid.UUID]SchoolShortInfo, len(schools))

	for _, school := range schools {
		mapOfSchools[school.ID] = school
	}

	mapOfOrganizations := make(map[uuid.UUID]EduOrganizationShortInfo, len(organizations))

	for _, organization := range organizations {
		mapOfOrganizations[organization.ID] = organization
	}

	for userIndex := 0; userIndex < len(u); userIndex++ {
		roleOrganizationID := u[userIndex].OrganizationID

		if roleOrganizationID != nil {
			if organization, exists := mapOfOrganizations[*u[userIndex].OrganizationID]; exists {
				u[userIndex].SetOrganizationShortInfo(organization)
			}
		}

		if u[userIndex].SchoolID == nil {
			continue
		}

		if school, exists := mapOfSchools[*u[userIndex].SchoolID]; exists {
			u[userIndex].SetSchoolShortInfo(school)

			if organization, exists := mapOfOrganizations[school.OrganizationID]; exists {
				u[userIndex].SetOrganizationShortInfo(organization)
			}
		}
	}
}

// NewUserRole creates new UserRole domain.
func NewUserRole(
	userID uuid.UUID,
	role Role,
	schoolID *uuid.UUID,
	organizationID *uuid.UUID,

	nowFunc func() time.Time,
) (UserRole, error) {
	now := nowFunc()

	if ok := role.Validate(); !ok {
		return UserRole{}, ErrUserRoleBadRequest
	}

	return UserRole{
		ID:             uuid.New(),
		UserID:         userID,
		Role:           role,
		SchoolID:       schoolID,
		OrganizationID: organizationID,

		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Role represents the user's role.
type Role string

// IsWithOnConflict checks if the role is with on conflict.
func (r Role) IsWithOnConflict() bool {
	switch r {
	case RoleGuardian:
		return true
	default:
		return false
	}
}

// Roles are list of Role.
type Roles []Role

// Validate validate user roles.
func (r Roles) Validate() bool {
	for _, role := range r {
		if !role.Validate() {
			return false
		}
	}

	return true
}

const (
	// RoleAdmin user role type admin.
	RoleAdmin Role = "admin"
	// RoleOwner user role type owner.
	RoleOwner Role = "owner"
	// RoleDirector user role type director.
	RoleDirector Role = "director"
	// RoleHeadmaster user role type headmaster.
	RoleHeadmaster Role = "headmaster"
	// RoleTeacher user role type teacher.
	RoleTeacher Role = "teacher"
	// RoleStudent user role type student.
	RoleStudent Role = "student"
	// RoleGuardian user role type guardian.
	RoleGuardian Role = "guardian"
)

// Validate validate user role.
func (r Role) Validate() bool {
	switch r {
	case RoleAdmin, RoleOwner, RoleDirector, RoleHeadmaster, RoleTeacher, RoleStudent, RoleGuardian:
		return true
	}

	return false
}

// UserListFilter filter for the list of User.
type UserListFilter struct {
	ListFilter

	SchoolIDs       []uuid.UUID
	OrganizationIDs []uuid.UUID
	Roles           Roles

	Emails []string
	Phones []string
}

// NewUserListFilter creates a new UserListFilter domain.
func NewUserListFilter(
	organizationIDs []uuid.UUID,
	schoolIDs []uuid.UUID,

	roles Roles,

	emails []string,
	phones []string,

	list ListFilter,
) (UserListFilter, error) {
	if !roles.Validate() {
		return UserListFilter{}, ErrUserRoleBadRequest
	}

	return UserListFilter{
		OrganizationIDs: organizationIDs,
		SchoolIDs:       schoolIDs,

		Roles: roles,

		Emails: emails,
		Phones: phones,

		ListFilter: list,
	}, nil
}

// IsStudent checks whether user is student.
func (u *UserRole) IsStudent() bool {
	return u.Role == RoleStudent
}
