package domain

import (
	"time"

	"github.com/google/uuid"
)

// Student is student domain.
type Student struct {
	ID              uuid.UUID
	RoleID          uuid.UUID
	UserID          uuid.UUID
	GroupID         uuid.UUID
	SchoolID        uuid.UUID
	SchoolShortInfo *SchoolShortInfo
	Group           Group

	User

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewStudent creates a new Student domain.
func NewStudent(
	roleID uuid.UUID,
	userID uuid.UUID,
	groupID uuid.UUID,
	nowFunc func() time.Time,
) Student {
	now := nowFunc()

	return Student{
		ID:        uuid.New(),
		RoleID:    roleID,
		UserID:    userID,
		GroupID:   groupID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetUser sets User info in Student model.
func (s *Student) SetUser(user User) {
	s.User = user
}

// SetGroup sets Group info in Student model.
func (s *Student) SetGroup(group Group) {
	s.Group = group
}

// SetShortSchool sets school info into student.
func (s *Student) SetShortSchool(schoolShort SchoolShortInfo) {
	s.SchoolShortInfo = &schoolShort
}

// Students are list of Student.
type Students []Student

// UserIDs returns the list of Student users ids.
func (s Students) UserIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(s))

	for _, student := range s {
		list = append(list, student.UserID)
	}

	return list
}

// GroupIDs returns the list of Groups users ids.
func (s Students) GroupIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(s))

	for _, student := range s {
		list = append(list, student.GroupID)
	}

	return list
}

// SetUsers sets users info into list of student.
func (s Students) SetUsers(users Users) {
	mapOfUsers := make(map[uuid.UUID]User, len(users))
	for _, user := range users {
		mapOfUsers[user.ID] = user
	}

	for index := range s {
		s[index].SetUser(mapOfUsers[s[index].UserID])
	}
}

// SetGroups sets group info into list of student.
func (s Students) SetGroups(groups Groups) {
	mapOfGroup := make(map[uuid.UUID]Group, len(groups))
	for _, group := range groups {
		mapOfGroup[group.ID] = group
	}

	for index := range s {
		s[index].SetGroup(mapOfGroup[s[index].GroupID])
	}
}

// SchoolIDs returns the list of student schools ids.
func (s Students) SchoolIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(s))

	for _, student := range s {
		list = append(list, student.SchoolID)
	}

	return list
}

// SetShortSchools sets school info into list of students.
func (s Students) SetShortSchools(schoolShorts SchoolShortInfos) {
	mapOfSchoolShort := make(map[uuid.UUID]SchoolShortInfo, len(schoolShorts))
	for _, schoolShort := range schoolShorts {
		mapOfSchoolShort[schoolShort.ID] = schoolShort
	}

	for index := range s {
		s[index].SetShortSchool(mapOfSchoolShort[s[index].SchoolID])
	}
}

// StudentListFilter filter for the list of Headmaster.
type StudentListFilter struct {
	CreatedDate DateFilter
	ListFilter
	GroupIDs        []uuid.UUID
	SchoolIDs       []uuid.UUID
	OrganizationIDs []uuid.UUID
}

// NewStudentListFilter creates a new StudentListFilter domain.
func NewStudentListFilter(
	createdDate DateFilter,
	list ListFilter,

	groupIDs []uuid.UUID,
	schoolIDs []uuid.UUID,
	organizationIDs []uuid.UUID,
) StudentListFilter {
	return StudentListFilter{
		CreatedDate: createdDate,
		ListFilter:  list,

		GroupIDs:        groupIDs,
		SchoolIDs:       schoolIDs,
		OrganizationIDs: organizationIDs,
	}
}
