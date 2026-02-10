package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// Subject is structure of Subject.
type Subject struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewSubject creates a new subject response from domain subject.
func NewSubject(subject domain.Subject) Subject {
	return Subject{
		ID:          subject.ID,
		Name:        subject.Name,
		Description: subject.Description,

		CreatedAt: utils.RFC3339Time(subject.CreatedAt),
		UpdatedAt: utils.RFC3339Time(subject.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(subject.DeletedAt),
	}
}

// SubjectList is a list of Subject.
type SubjectList struct {
	Subjects   []Subject  `json:"subjects"`
	Pagination Pagination `json:"pagination"`
}

// NewSubjectList creates a new subject list response from domain subject list.
func NewSubjectList(subjects []domain.Subject, responsePagination Pagination) SubjectList {
	responseSubjects := make([]Subject, 0, len(subjects))

	for _, subject := range subjects {
		responseSubjects = append(responseSubjects, NewSubject(subject))
	}

	return SubjectList{
		Subjects:   responseSubjects,
		Pagination: responsePagination,
	}
}
