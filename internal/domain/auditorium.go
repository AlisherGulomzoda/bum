package domain

import (
	"time"

	"github.com/google/uuid"
)

// Auditorium is structure of School auditorium.
type Auditorium struct {
	ID              uuid.UUID
	SchoolID        uuid.UUID
	Name            string
	SchoolSubjectID *uuid.UUID
	Description     *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Auditoriums is a list of Auditorium.
type Auditoriums []Auditorium

// AuditoriumListFilters is structure of Auditoriums list filters.
type AuditoriumListFilters struct {
	ListFilter
	SchoolID uuid.UUID
}

// NewAuditoriumListFilters creates a new AuditoriumFilters domain.
func NewAuditoriumListFilters(filter ListFilter, schoolID uuid.UUID) AuditoriumListFilters {
	return AuditoriumListFilters{
		ListFilter: filter,
		SchoolID:   schoolID,
	}
}

// NewAuditorium creates a new Auditorium.
func NewAuditorium(
	schoolID uuid.UUID,
	name string,
	schoolSubjectID *uuid.UUID,
	description *string,
	nowFunc func() time.Time,
) Auditorium {
	now := nowFunc()

	return Auditorium{
		ID:              uuid.New(),
		SchoolID:        schoolID,
		Name:            name,
		SchoolSubjectID: schoolSubjectID,
		Description:     description,

		CreatedAt: now,
		UpdatedAt: now,
	}
}
