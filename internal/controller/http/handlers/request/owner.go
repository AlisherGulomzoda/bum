package request

import (
	"github.com/google/uuid"

	"bum-service/pkg/liblog"
)

// AddOwner is a request to add a new owner.
type AddOwner struct {
	UserID         uuid.UUID `json:"user_id" binding:"required,uuid"`
	OrganizationID uuid.UUID `json:"organization_id" binding:"required,uuid"`
	Phone          *string   `json:"phone" binding:"omitempty,e164"`
	Email          *string   `json:"email" binding:"omitempty,email"`
}

// LogFields returns a list of fields for logging.
func (c AddOwner) LogFields() liblog.Fields {
	return liblog.Fields{
		"user_id":         c.UserID,
		"organization_id": c.OrganizationID,
		"phone":           c.Phone,
		"email":           c.Email,
	}
}

// OwnerList request model for listing of owners.
type OwnerList struct {
	ListFilter

	CreatedDate DateFilter

	OrganizationIDsFilter
}
