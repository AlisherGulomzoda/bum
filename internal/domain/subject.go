package domain

import (
	"time"

	"github.com/google/uuid"
)

// Subject is a domain model.
type Subject struct {
	ID          uuid.UUID
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// NewSubject creates a new Subject domain.
func NewSubject(
	name string,
	description *string,

	nowFunc func() time.Time,
) Subject {
	now := nowFunc()

	return Subject{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Subjects are collections of Subject.
type Subjects []Subject

// SubjectListFilter filter for the list of Subject.
type SubjectListFilter struct {
	ListFilter
}

// NewSubjectListFilter creates a new SubjectListFilter domain.
func NewSubjectListFilter(list ListFilter) SubjectListFilter {
	return SubjectListFilter{ListFilter: list}
}
