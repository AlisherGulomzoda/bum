package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// EduOrganization is structure of educational organization response.
type EduOrganization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Logo        *string   `json:"logo,omitempty"`
	Description *string   `json:"description,omitempty"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewEduOrganization creates a new educational organization response.
func NewEduOrganization(
	organization domain.EduOrganization,
) EduOrganization {
	return EduOrganization{
		ID:          organization.ID,
		Name:        organization.Name,
		Logo:        organization.Logo,
		Description: organization.Description,

		CreatedAt: utils.RFC3339Time(organization.CreatedAt),
		UpdatedAt: utils.RFC3339Time(organization.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(organization.DeletedAt),
	}
}

// EduOrganizations is a collection of Organizations.
type EduOrganizations []EduOrganization

// EduOrganizationsList struct of Organizations lists.
type EduOrganizationsList struct {
	EduOrganizations EduOrganizations `json:"organizations"`
	Pagination       Pagination       `json:"pagination"`
}

// NewEduOrganizationsList returns a new EduOrganizations list response.
func NewEduOrganizationsList(
	domains domain.EduOrganizations,
	pagination Pagination,
) EduOrganizationsList {
	organizations := make(EduOrganizations, 0, len(domains))

	for _, org := range domains {
		organizations = append(organizations, NewEduOrganization(org))
	}

	return EduOrganizationsList{
		EduOrganizations: organizations,
		Pagination:       pagination,
	}
}

// EduOrganizationShortInfo short information about organization.
type EduOrganizationShortInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NewEduOrganizationShortInfo creates a new EduOrganizationShortInfo from domain.
func NewEduOrganizationShortInfo(organization *domain.EduOrganizationShortInfo) *EduOrganizationShortInfo {
	if organization == nil {
		return nil
	}

	return &EduOrganizationShortInfo{
		ID:   organization.ID,
		Name: organization.Name,
	}
}
