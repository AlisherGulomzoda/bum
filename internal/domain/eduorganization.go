package domain

import (
	"time"

	"github.com/google/uuid"
)

// EduOrganization is structure of educational organization.
type EduOrganization struct {
	ID          uuid.UUID
	Name        string
	Logo        *string
	Description *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewEduOrganization creates a new educational organization.
func NewEduOrganization(
	name string,
	logo *string,
	description *string,

	nowFunc func() time.Time,
) EduOrganization {
	now := nowFunc()

	return EduOrganization{
		ID:          uuid.New(),
		Name:        name,
		Logo:        logo,
		Description: description,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// EduOrganizations is a list of EduOrganization.
type EduOrganizations []EduOrganization

// EduOrganizationFilters is structure of EduOrganization filters.
type EduOrganizationFilters struct {
	ListFilter
}

// NewEduOrganizationFilters creates a new EduOrganizationFilters domain.
func NewEduOrganizationFilters(list ListFilter) EduOrganizationFilters {
	return EduOrganizationFilters{ListFilter: list}
}

// Update updates organization fields.
func (e *EduOrganization) Update(
	name string,
	logo *string,
	nowFunc func() time.Time,
) {
	now := nowFunc()

	e.Name = name
	e.Logo = logo

	e.UpdatedAt = now
}

// EduOrganizationShortInfo is structure of educational organization short information.
type EduOrganizationShortInfo struct {
	ID   uuid.UUID
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// EduOrganizationShortInfos is list of EduOrganizationShortInfo.
type EduOrganizationShortInfos []EduOrganizationShortInfo
