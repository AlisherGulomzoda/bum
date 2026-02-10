//nolint:dupl // it's ok
package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// Director is a structure of director response.
type Director struct {
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

// NewDirector creates a new director response.
func NewDirector(director domain.Director) Director {
	return Director{
		ID:              director.ID,
		RoleID:          director.RoleID,
		User:            NewUser(director.User),
		SchoolID:        director.SchoolID,
		SchoolShortInfo: NewSchoolShortInfo(director.SchoolShortInfo),
		Phone:           director.Phone,
		Email:           director.Email,

		CreatedAt: utils.RFC3339Time(director.CreatedAt),
		UpdatedAt: utils.RFC3339Time(director.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(director.DeletedAt),
	}
}

// DirectorList response model for listing directors.
type DirectorList struct {
	Directors  []Director `json:"directors"`
	Pagination Pagination `json:"pagination"`
}

// NewDirectorList creates a new director list for response.
func NewDirectorList(
	list []domain.Director,
	pagination Pagination,
) DirectorList {
	directors := make([]Director, len(list))

	for i := range list {
		directors[i] = NewDirector(list[i])
	}

	return DirectorList{
		Directors:  directors,
		Pagination: pagination,
	}
}
