package request

import (
	"github.com/google/uuid"

	"bum-service/pkg/liblog"
)

// CreateEduOrganizational is a request to create a new educational organization.
type CreateEduOrganizational struct {
	Name        string  `json:"name" binding:"required"`
	Logo        *string `json:"logo" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
}

// LogFields returns a list of fields for logging.
func (c CreateEduOrganizational) LogFields() liblog.Fields {
	return liblog.Fields{
		"name":        c.Name,
		"logo":        c.Logo,
		"description": c.Description,
	}
}

// GetEduOrganizationalList is a request to get a list of educational organizations.
type GetEduOrganizationalList struct {
	ListFilter
}

// UpdateEduOrganizational is a request to update educational organization.
type UpdateEduOrganizational struct {
	Name string  `json:"name" binding:"required"`
	Logo *string `json:"logo,omitempty"`
}

// OrganizationIDsFilter is a filter for organization IDs.
type OrganizationIDsFilter struct {
	OrganizationIDs []string `form:"organization_ids[]" binding:"omitempty,dive,uuid"`
}

// OrganizationUUIDs converts organizationID to list of uuid.
func (s OrganizationIDsFilter) OrganizationUUIDs() []uuid.UUID {
	if len(s.OrganizationIDs) == 0 {
		return []uuid.UUID{}
	}

	uuids := make([]uuid.UUID, len(s.OrganizationIDs))
	for idx := range s.OrganizationIDs {
		uuids[idx] = uuid.MustParse(s.OrganizationIDs[idx])
	}

	return uuids
}
