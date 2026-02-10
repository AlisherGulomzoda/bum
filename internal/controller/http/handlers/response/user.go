package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// User is a structure of user response.
type User struct {
	ID         uuid.UUID  `json:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	MiddleName *string    `json:"middle_name,omitempty"`
	Gender     string     `json:"gender"`
	Phone      *string    `json:"phone,omitempty"`
	Email      string     `json:"email"`
	UserRoles  []UserRole `json:"user_roles,omitempty"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewUser creates a new user response.
func NewUser(user domain.User) User {
	return User{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Gender:     string(user.Gender),
		Phone:      user.Phone,
		Email:      user.Email,
		UserRoles:  NewUserRoles(user.UserRoles),

		CreatedAt: utils.RFC3339Time(user.CreatedAt),
		UpdatedAt: utils.RFC3339Time(user.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(user.DeletedAt),
	}
}

// UserList is a list of User.
type UserList struct {
	Users      []User     `json:"users"`
	Pagination Pagination `json:"pagination"`
}

// NewUserList creates a new user list response from domain user list.
func NewUserList(users domain.Users, responsePagination Pagination) UserList {
	responseUsers := make([]User, 0, len(users))

	for id := range users {
		responseUsers = append(responseUsers, NewUser(users[id]))
	}

	return UserList{
		Users:      responseUsers,
		Pagination: responsePagination,
	}
}

// GroupShortInfo is a structure of group short info.
type GroupShortInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// UserRole is list of user roles response.
type UserRole struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`

	SchoolID        *uuid.UUID       `json:"school_id,omitempty"`
	SchoolShortInfo *SchoolShortInfo `json:"school_short_info,omitempty"`

	OrganizationID        *uuid.UUID                `json:"organization_id,omitempty"`
	OrganizationShortInfo *EduOrganizationShortInfo `json:"organization_short_info,omitempty"`

	GroupID *uuid.UUID      `json:"group_id,omitempty"`
	Group   *GroupShortInfo `json:"group,omitempty"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewUserRoles creates a new user roles response from domain user roles.
func NewUserRoles(userRoles domain.UserRoles) []UserRole {
	responseUserRoles := make([]UserRole, 0, len(userRoles))

	for id := range userRoles {
		role := UserRole{
			ID:     userRoles[id].ID,
			UserID: userRoles[id].UserID,
			Role:   string(userRoles[id].Role),

			SchoolID:        userRoles[id].SchoolID,
			SchoolShortInfo: NewSchoolShortInfo(userRoles[id].SchoolShortInfo),

			OrganizationID:        userRoles[id].OrganizationID,
			OrganizationShortInfo: NewEduOrganizationShortInfo(userRoles[id].EduOrganizationShortInfo),

			GroupID: userRoles[id].GroupID,

			CreatedAt: utils.RFC3339Time(userRoles[id].CreatedAt),
			UpdatedAt: utils.RFC3339Time(userRoles[id].UpdatedAt),
			DeletedAt: (*utils.RFC3339Time)(userRoles[id].DeletedAt),
		}

		if userRoles[id].GroupID != nil {
			group := GroupShortInfo{
				ID:   userRoles[id].Group.ID,
				Name: userRoles[id].Group.Name,

				CreatedAt: utils.RFC3339Time(userRoles[id].Group.CreatedAt),
				UpdatedAt: utils.RFC3339Time(userRoles[id].Group.UpdatedAt),
				DeletedAt: (*utils.RFC3339Time)(userRoles[id].Group.DeletedAt),
			}

			role.Group = &group
		}

		responseUserRoles = append(responseUserRoles,
			role,
		)
	}

	return responseUserRoles
}
