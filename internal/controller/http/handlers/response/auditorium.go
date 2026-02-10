package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// Auditorium is a school auditorium object response.
type Auditorium struct {
	ID              uuid.UUID  `json:"id"`
	SchoolID        uuid.UUID  `json:"school_id"`
	Name            string     `json:"name"`
	SchoolSubjectID *uuid.UUID `json:"school_subject_id_id,omitempty"`
	Description     *string    `json:"description,omitempty"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewSchoolAuditorium creates a new school auditorium response.
func NewSchoolAuditorium(auditorium domain.Auditorium) Auditorium {
	rep := Auditorium{
		ID:              auditorium.ID,
		SchoolID:        auditorium.SchoolID,
		Name:            auditorium.Name,
		SchoolSubjectID: auditorium.SchoolSubjectID,
		Description:     auditorium.Description,

		CreatedAt: utils.RFC3339Time(auditorium.CreatedAt),
		UpdatedAt: utils.RFC3339Time(auditorium.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(auditorium.DeletedAt),
	}

	return rep
}

// SchoolAuditoriumList is response of school auditorium list.
type SchoolAuditoriumList struct {
	Auditoriums []Auditorium `json:"auditoriums"`
	Pagination  Pagination   `json:"pagination"`
}

// NewSchoolAuditoriums returns a new set of Auditoriums response.
func NewSchoolAuditoriums(auditoriums domain.Auditoriums) []Auditorium {
	list := make([]Auditorium, 0, len(auditoriums))

	for _, auditorium := range auditoriums {
		list = append(list, NewSchoolAuditorium(auditorium))
	}

	return list
}

// NewSchoolAuditoriumList creates a new SchoolAuditoriumList response.
func NewSchoolAuditoriumList(list domain.Auditoriums, pagination Pagination) SchoolAuditoriumList {
	return SchoolAuditoriumList{
		Auditoriums: NewSchoolAuditoriums(list),
		Pagination:  pagination,
	}
}
