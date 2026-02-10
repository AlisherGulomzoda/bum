package request

import (
	"github.com/google/uuid"

	"bum-service/pkg/liblog"
)

// CreateGroup is a request to create a new group.
type CreateGroup struct {
	Name    string    `json:"name" binding:"required"`
	GradeID uuid.UUID `json:"grade_id" binding:"required"`
}

// LogFields returns a list of fields for logging.
func (c CreateGroup) LogFields() liblog.Fields {
	return liblog.Fields{
		"name":     c.Name,
		"grade_id": c.GradeID,
	}
}

// AddGroupSubject is a request to add subject to a group.
type AddGroupSubject struct {
	SchoolSubjectID uuid.UUID  `json:"school_subject_id" binding:"required,uuid"`
	TeacherID       *uuid.UUID `json:"teacher_id,omitempty" binding:"omitnil,uuid"`
	Count           *int16     `json:"count,omitempty" binding:"omitnil"`
}

// UpdateGroup updates group.
type UpdateGroup struct {
	Name    string    `json:"name" binding:"required"`
	GradeID uuid.UUID `json:"grade_id" binding:"required"`

	ClassTeacherID         *uuid.UUID `json:"class_teacher_id" binding:"omitnil,uuid"`
	ClassPresidentID       *uuid.UUID `json:"class_president_id" binding:"omitnil,uuid"`
	DeputyClassPresidentID *uuid.UUID `json:"deputy_class_president_id" binding:"omitnil,uuid"`
}

// GroupIDsFilter group ids filter for embedding.
type GroupIDsFilter struct {
	GroupIDs []string `form:"group_ids[]" binding:"omitempty,dive,uuid"`
}

// GroupUUIDs converts GroupIDs to list of uuid.
func (g GroupIDsFilter) GroupUUIDs() []uuid.UUID { // todo: для удобства можно юзать дженерики.
	if len(g.GroupIDs) == 0 {
		return []uuid.UUID{}
	}

	uuids := make([]uuid.UUID, len(g.GroupIDs))
	for idx := range g.GroupIDs {
		uuids[idx] = uuid.MustParse(g.GroupIDs[idx])
	}

	return uuids
}
