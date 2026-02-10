package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// Owner is a structure of owner response.
type Owner struct {
	ID             uuid.UUID       `json:"id"`
	RoleID         uuid.UUID       `json:"role_id"`
	User           User            `json:"user"`
	Phone          *string         `json:"phone,omitempty"`
	Email          *string         `json:"email,omitempty"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	Organization   EduOrganization `json:"organization"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewOwner creates a new owner response.
func NewOwner(owner domain.Owner) Owner {
	return Owner{
		ID:             owner.ID,
		RoleID:         owner.RoleID,
		User:           NewUser(owner.User),
		Phone:          owner.Phone,
		Email:          owner.Email,
		OrganizationID: owner.OrganizationID,

		CreatedAt: utils.RFC3339Time(owner.CreatedAt),
		UpdatedAt: utils.RFC3339Time(owner.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(owner.DeletedAt),
	}
}

// OwnerList response model for listing owners.
type OwnerList struct {
	Owners     []Owner    `json:"owners"`
	Pagination Pagination `json:"pagination"`
}

// NewOwnerList creates a new owner list for response.
func NewOwnerList(
	list domain.Owners,
	pagination Pagination,
) OwnerList {
	owners := make([]Owner, len(list))

	for i := range list {
		owners[i] = NewOwner(list[i])
	}

	return OwnerList{
		Owners:     owners,
		Pagination: pagination,
	}
}
