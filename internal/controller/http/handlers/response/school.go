package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// School is structure of School.
type School struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Location        string     `json:"location"`
	Phone           *string    `json:"phone,omitempty"`
	Email           *string    `json:"email,omitempty"`
	GradeStandardID *uuid.UUID `json:"grade_standard_id"`

	OrganizationID uuid.UUID                `json:"organization_id"`
	Organization   EduOrganizationShortInfo `json:"organization"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewSchool creates a new school response from domain school.
func NewSchool(school domain.School) School {
	return School{
		ID:              school.ID,
		Name:            school.Name,
		OrganizationID:  school.OrganizationID,
		Organization:    *NewEduOrganizationShortInfo(school.OrganizationShortInfo),
		Location:        school.Location,
		Phone:           school.Phone,
		Email:           school.Email,
		GradeStandardID: school.GradeStandardID,

		CreatedAt: utils.RFC3339Time(school.CreatedAt),
		UpdatedAt: utils.RFC3339Time(school.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(school.DeletedAt),
	}
}

// SchoolList response model for listing schools.
type SchoolList struct {
	Schools    []School   `json:"schools"`
	Pagination Pagination `json:"pagination"`
}

// NewSchools creates a new school response from domain school.
func NewSchools(
	schools domain.Schools,
	pagination Pagination,
) SchoolList {
	schoolList := make([]School, 0, len(schools))

	for i := range schools {
		schoolList = append(schoolList, NewSchool(schools[i]))
	}

	return SchoolList{
		Schools:    schoolList,
		Pagination: pagination,
	}
}

// SchoolShortInfo short information about school.
type SchoolShortInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NewSchoolShortInfo creates a new SchoolShortInfo from domain.
func NewSchoolShortInfo(school *domain.SchoolShortInfo) *SchoolShortInfo {
	if school == nil {
		return nil
	}

	return &SchoolShortInfo{
		ID:   school.ID,
		Name: school.Name,
	}
}

// SchoolSubject is a SchoolSubject object response.
type SchoolSubject struct {
	ID          uuid.UUID `json:"id"`
	SubjectID   uuid.UUID `json:"subject_id"`
	SchoolID    uuid.UUID `json:"school_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewSchoolSubject creates a new schoolSubject resource.
func NewSchoolSubject(schoolSubject domain.SchoolSubject) SchoolSubject {
	return SchoolSubject{
		ID:          schoolSubject.ID,
		SubjectID:   schoolSubject.SubjectID,
		SchoolID:    schoolSubject.SchoolID,
		Name:        schoolSubject.Name,
		Description: schoolSubject.Description,

		CreatedAt: utils.RFC3339Time(schoolSubject.CreatedAt),
		UpdatedAt: utils.RFC3339Time(schoolSubject.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(schoolSubject.DeletedAt),
	}
}

// NewSchoolSubjects returns a new set of SchoolSubjects response.
func NewSchoolSubjects(schoolSubjects domain.SchoolSubjects) []SchoolSubject {
	list := make([]SchoolSubject, 0, len(schoolSubjects))

	for _, schoolSubject := range schoolSubjects {
		list = append(list, NewSchoolSubject(schoolSubject))
	}

	return list
}

// SchoolSubjectList is response of school subjects list.
type SchoolSubjectList struct {
	SchoolSubjects []SchoolSubject `json:"school_subjects"`
	Pagination     Pagination      `json:"pagination"`
}

// NewSchoolSubjectsList creates a new SchoolSubjectList resource.
func NewSchoolSubjectsList(schoolSubjects domain.SchoolSubjects, pagination Pagination) SchoolSubjectList {
	return SchoolSubjectList{
		SchoolSubjects: NewSchoolSubjects(schoolSubjects),
		Pagination:     pagination,
	}
}

// StudyPlan is a response model for study plan.
type StudyPlan struct {
	ID           string  `json:"id"`
	GroupSubject string  `json:"group_subject"`
	Title        string  `json:"title"`
	Description  *string `json:"description,omitempty"`
	PlanOrder    int16   `json:"plan_order"`
	Status       string  `json:"status"`
}

// NewStudyPlans creates a new StudyPlan response from domain study plans.
func NewStudyPlans(studyPlans domain.StudyPlans) []StudyPlan {
	resp := make([]StudyPlan, len(studyPlans))

	for i, plan := range studyPlans {
		resp[i] = StudyPlan{
			ID:           plan.ID.String(),
			GroupSubject: plan.GroupSubjectID.String(),
			Title:        plan.Title,
			Description:  plan.Description,
			PlanOrder:    plan.PlanOrder,
			Status:       plan.Status.String(),
		}
	}

	return resp
}

// NewStudyPlan creates a new StudyPlan response from domain study plan.
func NewStudyPlan(studyPlan domain.StudyPlan) StudyPlan {
	return StudyPlan{
		ID:           studyPlan.ID.String(),
		GroupSubject: studyPlan.GroupSubjectID.String(),
		Title:        studyPlan.Title,
		Description:  studyPlan.Description,
		PlanOrder:    studyPlan.PlanOrder,
		Status:       studyPlan.Status.String(),
	}
}
