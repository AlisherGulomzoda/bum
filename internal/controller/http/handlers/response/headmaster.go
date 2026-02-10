//nolint:dupl // it's ok
package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// Headmaster is a structure of headmaster response.
type Headmaster struct {
	ID              uuid.UUID        `json:"id"`
	RoleID          uuid.UUID        `json:"role_id"`
	User            User             `json:"user"`
	SchoolID        uuid.UUID        `json:"school_id"`
	SchoolShortInfo *SchoolShortInfo `json:"school_short_info,omitempty"`
	Phone           *string          `json:"phone,omitempty"`
	Email           *string          `json:"email,omitempty"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewHeadmaster creates a new headmaster response.
func NewHeadmaster(headmaster domain.Headmaster) Headmaster {
	return Headmaster{
		ID:              headmaster.ID,
		RoleID:          headmaster.RoleID,
		User:            NewUser(headmaster.User),
		SchoolID:        headmaster.SchoolID,
		SchoolShortInfo: NewSchoolShortInfo(headmaster.SchoolShortInfo),
		Phone:           headmaster.Phone,
		Email:           headmaster.Email,

		CreatedAt: utils.RFC3339Time(headmaster.CreatedAt),
		UpdatedAt: utils.RFC3339Time(headmaster.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(headmaster.DeletedAt),
	}
}

// HeadmasterList response model for listing headmasters.
type HeadmasterList struct {
	Headmasters []Headmaster `json:"headmasters"`
	Pagination  Pagination   `json:"pagination"`
}

// NewHeadmasterList creates a new headmaster list for response.
func NewHeadmasterList(
	list []domain.Headmaster,
	pagination Pagination,
) HeadmasterList {
	headmasters := make([]Headmaster, len(list))

	for i := range list {
		headmasters[i] = NewHeadmaster(list[i])
	}

	return HeadmasterList{
		Headmasters: headmasters,
		Pagination:  pagination,
	}
}
