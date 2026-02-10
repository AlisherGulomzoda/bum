package request

import (
	"bum-service/internal/domain"
	"bum-service/pkg/liblog"
)

// AddUser  is a request to create a new user.
type AddUser struct {
	FirstName  string  `json:"first_name" binding:"required"`
	LastName   string  `json:"last_name" binding:"required"`
	MiddleName *string `json:"middle_name" binding:"omitempty"`
	Gender     string  `json:"gender" binding:"required"`
	Phone      *string `json:"phone" binding:"omitempty,e164"`
	Email      string  `json:"email" binding:"email,required"`
	Password   string  `json:"password" binding:"required"`
}

// LogInfo returns user information for logging.
func (c AddUser) LogInfo() liblog.Fields {
	return liblog.Fields{
		"first_name":  c.FirstName,
		"last_name":   c.LastName,
		"middle_name": c.MiddleName,
		"gender":      c.Gender,
		"phone":       c.Phone,
		"email":       c.Email,
	}
}

// UserList is a request for User list.
type UserList struct {
	ListFilter

	OrganizationIDsFilter
	SchoolIDsFilter
	Roles []string `form:"roles[]" binding:"omitempty"`

	Emails []string `form:"emails[]" binding:"omitempty,dive,email"`
	Phones []string `form:"phones[]" binding:"omitempty,dive,e164"`
}

// UserRoles returns list fo user roles for filter.
func (u UserList) UserRoles() domain.Roles {
	roles := make(domain.Roles, 0, len(u.Roles))

	for _, role := range u.Roles {
		roles = append(roles, domain.Role(role))
	}

	return roles
}
