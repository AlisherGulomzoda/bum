package domain

import (
	"time"

	"github.com/google/uuid"
)

// School is structure of School.
type School struct {
	ID              uuid.UUID
	Name            string
	GradeStandardID *uuid.UUID

	OrganizationID        uuid.UUID
	OrganizationShortInfo *EduOrganizationShortInfo

	Location string
	Phone    *string
	Email    *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// SetOrganizationShortInfo sets organization short info.
func (s *School) SetOrganizationShortInfo(organization EduOrganizationShortInfo) {
	s.OrganizationShortInfo = &organization
}

// Schools are collection of School.
type Schools []School

// GetOrganizationIDs returns list of organizations ids.
func (s Schools) GetOrganizationIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(s))

	for _, school := range s {
		list = append(list, school.OrganizationID)
	}

	return list
}

// SetOrganizationShortInfos sets organizations short info to list of Schools.
func (s Schools) SetOrganizationShortInfos(organizations EduOrganizationShortInfos) {
	mapOfOrganizations := make(map[uuid.UUID]EduOrganizationShortInfo)
	for _, org := range organizations {
		mapOfOrganizations[org.ID] = org
	}

	for index := range s {
		s[index].SetOrganizationShortInfo(mapOfOrganizations[s[index].OrganizationID])
	}
}

// SchoolFilterParams is structure of School filter params.
type SchoolFilterParams struct {
	Pagination
}

// NewSchool creates a new School domain.
func NewSchool(
	name string,
	organizationID uuid.UUID,
	location string,
	phone *string,
	email *string,
	nowFunc func() time.Time,
) School {
	now := nowFunc()

	return School{
		ID:             uuid.New(),
		Name:           name,
		OrganizationID: organizationID,
		Location:       location,
		Phone:          phone,
		Email:          email,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates school domain.
func (s *School) Update(
	name string,
	location string,
	phone *string,
	email *string,
	gradeStandardID *uuid.UUID,

	nowFunc func() time.Time,
) {
	now := nowFunc()

	s.Name = name
	s.Location = location
	s.Phone = phone
	s.Email = email
	s.GradeStandardID = gradeStandardID

	s.UpdatedAt = now
}

// SchoolFilters is structure of School filters.
type SchoolFilters struct {
	ListFilter

	Emails          []string
	Phones          []string
	OrganizationIDs []uuid.UUID
}

// NewSchoolFilters creates a new SchoolFilters domain.
func NewSchoolFilters(filter ListFilter, emails, phones []string, organizationIDs []uuid.UUID) SchoolFilters {
	return SchoolFilters{
		ListFilter:      filter,
		Emails:          emails,
		Phones:          phones,
		OrganizationIDs: organizationIDs,
	}
}

// SchoolShortInfo is structure of short School info.
type SchoolShortInfo struct {
	ID             uuid.UUID
	Name           string
	OrganizationID uuid.UUID
}

// SchoolShortInfos is list of SchoolShortInfo.
type SchoolShortInfos []SchoolShortInfo

// OrganizationIDs returns list of organization ids.
func (s SchoolShortInfos) OrganizationIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(s))

	for _, school := range s {
		list = append(list, school.OrganizationID)
	}

	return list
}

// SchoolSubject is structure of School subject.
type SchoolSubject struct {
	ID          uuid.UUID
	SchoolID    uuid.UUID
	SubjectID   uuid.UUID
	Name        string
	Description *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewSchoolSubject creates a new SchoolSubject.
func NewSchoolSubject(
	schoolID uuid.UUID,
	subjectID uuid.UUID,
	name string,
	description *string,
	nowFunc func() time.Time,
) SchoolSubject {
	now := nowFunc()

	return SchoolSubject{
		ID:          uuid.New(),
		SchoolID:    schoolID,
		SubjectID:   subjectID,
		Name:        name,
		Description: description,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SchoolSubjects is a collection of SchoolSubject.
type SchoolSubjects []SchoolSubject

// SchoolSubjectFilters is structure of School subject filters.
type SchoolSubjectFilters struct {
	ListFilter

	SchoolID uuid.UUID
}

// NewSchoolSubjectFilters creates a new SchoolSubjectFilters domain.
func NewSchoolSubjectFilters(filter ListFilter, schoolID uuid.UUID) SchoolSubjectFilters {
	return SchoolSubjectFilters{
		ListFilter: filter,
		SchoolID:   schoolID,
	}
}

// StudyPlan is a study plan domain model.
type StudyPlan struct {
	ID             uuid.UUID
	GroupSubjectID uuid.UUID
	Title          string
	Description    *string
	PlanOrder      int16
	Status         StudyPlanStatus

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// StudyPlanStatus is a study plan status domain.
type StudyPlanStatus string

// Validate validates StudyPlanStatus.
func (s StudyPlanStatus) Validate() error {
	switch s {
	case Planned, Ongoing, Completed:
		return nil
	default:
		return ErrInvalidStudyPlanStatus
	}
}

// String returns string representation of StudyPlanStatus.
func (s StudyPlanStatus) String() string {
	return string(s)
}

const (
	Planned   StudyPlanStatus = "planned"   // Planned is a status of planned study plan.
	Ongoing   StudyPlanStatus = "ongoing"   // Ongoing is a status of ongoing study plan.
	Completed StudyPlanStatus = "completed" // Completed is a status of completed study plan.
)

// NewStudyPlan creates a new study plan.
func NewStudyPlan(
	id *uuid.UUID,
	groupSubjectID uuid.UUID,
	title string,
	desc *string,
	order int16,
	nowFunc func() time.Time,
) StudyPlan {
	now := nowFunc()

	newStudyPlan := StudyPlan{
		ID:             uuid.New(),
		GroupSubjectID: groupSubjectID,
		Title:          title,
		Description:    desc,
		PlanOrder:      order,
		Status:         Planned, // default status is planned
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if id != nil {
		newStudyPlan.ID = *id
	}

	return newStudyPlan
}

// StudyPlans is a collection of StudyPlan.
type StudyPlans []StudyPlan
