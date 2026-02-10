package response

import (
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// Group is structure of Group.
type Group struct {
	ID       uuid.UUID `json:"id"`
	SchoolID uuid.UUID `json:"school_id"`
	Name     string    `json:"name"`
	GradeID  uuid.UUID `json:"grade_id"`

	ClassTeacherID *uuid.UUID        `json:"class_teacher_id"`
	ClassTeacher   *TeacherShortInfo `json:"class_teacher,omitempty"`

	ClassPresidentID *uuid.UUID        `json:"class_president_id"`
	ClassPresident   *StudentShortInfo `json:"class_president,omitempty"`

	DeputyClassPresidentID *uuid.UUID        `json:"deputy_class_president_id"`
	DeputyClassPresident   *StudentShortInfo `json:"deputy_class_president,omitempty"`

	Grade Grade `json:"grade"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewGroup creates a new group response from domain group.
func NewGroup(group domain.Group) Group {
	return Group{
		ID:       group.ID,
		SchoolID: group.SchoolID,
		Name:     group.Name,
		GradeID:  group.GradeID,

		ClassTeacherID: group.ClassTeacherID,
		ClassTeacher:   NewTeacherShortInfo(group.ClassTeacher),

		ClassPresidentID: group.ClassPresidentID,
		ClassPresident:   NewStudentShortInfo(group.ClassPresident),

		DeputyClassPresidentID: group.DeputyClassPresidentID,
		DeputyClassPresident:   NewStudentShortInfo(group.DeputyClassPresident),

		Grade: NewGrade(group.Grade),

		CreatedAt: utils.RFC3339Time(group.CreatedAt),
		UpdatedAt: utils.RFC3339Time(group.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(group.DeletedAt),
	}
}

// Groups is list of Group.
type Groups []Group

// NewGroups is constructor for creating groups.
func NewGroups(groups domain.Groups) Groups {
	list := make(Groups, 0, len(groups))

	for _, g := range groups {
		list = append(list, NewGroup(g))
	}

	return list
}

// GroupSubject is structure of GroupSubject.
type GroupSubject struct {
	ID      uuid.UUID `json:"id"`
	GroupID uuid.UUID `json:"group_id"`
	Count   *int16    `json:"count,omitempty"`

	SchoolSubjectID uuid.UUID     `json:"school_subject_id"`
	SchoolSubject   SchoolSubject `json:"school_subject"`

	TeacherID        *uuid.UUID        `json:"teacher_id,omitempty"`
	TeacherShortInfo *TeacherShortInfo `json:"teacher_short_info"`

	CreatedAt utils.RFC3339Time  `json:"created_at"`
	UpdatedAt utils.RFC3339Time  `json:"updated_at"`
	DeletedAt *utils.RFC3339Time `json:"deleted_at,omitempty"`
}

// NewGroupSubject creates a new group subject response.
func NewGroupSubject(groupSubject domain.GroupSubject) GroupSubject {
	return GroupSubject{
		ID:      groupSubject.ID,
		GroupID: groupSubject.GroupID,
		Count:   groupSubject.Count,

		SchoolSubjectID: groupSubject.SchoolSubjectID,
		SchoolSubject:   NewSchoolSubject(groupSubject.SchoolSubject),

		TeacherID:        groupSubject.TeacherID,
		TeacherShortInfo: NewTeacherShortInfo(groupSubject.Teacher),

		CreatedAt: utils.RFC3339Time(groupSubject.CreatedAt),
		UpdatedAt: utils.RFC3339Time(groupSubject.UpdatedAt),
		DeletedAt: (*utils.RFC3339Time)(groupSubject.DeletedAt),
	}
}

// GroupSubjectList is list of GroupSubject.
type GroupSubjectList struct {
	GroupSubject []GroupSubject `json:"group_subjects"`
	Total        int            `json:"total"`
}

// NewGroupSubjects creates a new group subjects response from domain group subjects.
func NewGroupSubjects(groupSubjects domain.GroupSubjects) GroupSubjectList {
	list := make([]GroupSubject, 0, len(groupSubjects))

	for _, gs := range groupSubjects {
		list = append(list, NewGroupSubject(gs))
	}

	return GroupSubjectList{
		GroupSubject: list,
		Total:        len(list),
	}
}
