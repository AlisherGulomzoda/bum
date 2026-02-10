package request

import "github.com/google/uuid"

// AddStudent is add student request.
type AddStudent struct {
	UserID   uuid.UUID `json:"user_id" binding:"required,uuid"`
	SchoolID uuid.UUID `json:"school_id" binding:"required,uuid"`
	GroupID  uuid.UUID `json:"group_id" binding:"required,uuid"`
}

// StudentList request model for listing of students.
type StudentList struct {
	ListFilter

	GroupIDsFilter
	SchoolIDsFilter
	OrganizationIDsFilter

	CreatedDate DateFilter
}
