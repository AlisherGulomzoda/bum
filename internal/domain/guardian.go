package domain

import (
	"time"

	"github.com/google/uuid"
)

// StudentGuardian is student guardian domain.
type StudentGuardian struct {
	ID uuid.UUID

	StudentID uuid.UUID
	Student   Student

	UserID uuid.UUID
	User   User

	Relation StudentGuardianRelation

	SchoolID uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewStudentGuardian creates a new StudentGuardian domain.
func NewStudentGuardian(
	studentID uuid.UUID,
	userID uuid.UUID,
	relation StudentGuardianRelation,
	schoolID uuid.UUID,
	nowFunc func() time.Time,
) (StudentGuardian, error) {
	now := nowFunc()

	if ok := relation.Validate(); !ok {
		return StudentGuardian{}, ErrStudentGuardianRelationBadRequest
	}

	return StudentGuardian{
		ID:        uuid.New(),
		UserID:    userID,
		StudentID: studentID,
		Relation:  relation,
		SchoolID:  schoolID,

		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// StudentGuardians are list of StudentGuardian.
type StudentGuardians []StudentGuardian

// Guardian is guardian domain.
type Guardian struct {
	User
	StudentGuardians StudentGuardians
}

// SetUser sets user info.
func (g *Guardian) SetUser(user User) {
	g.User = user
}

// SetStudents sets guardian students.
func (g *Guardian) SetStudents(studentGuardians StudentGuardians) {
	g.StudentGuardians = studentGuardians
}

// StudentIDs returns the list of StudentGuardians student ids.
func (s StudentGuardians) StudentIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(s))

	for _, studentGuardian := range s {
		list = append(list, studentGuardian.StudentID)
	}

	return list
}

// SetStudent sets Student info in StudentGuardian model.
func (s *StudentGuardian) SetStudent(student Student) {
	s.Student = student
}

// SetStudents sets student info into list of StudentGuardians.
func (s StudentGuardians) SetStudents(students Students) {
	mapOfStudent := make(map[uuid.UUID]Student, len(students))
	for index := range students {
		mapOfStudent[students[index].ID] = students[index]
	}

	for index := range s {
		s[index].SetStudent(mapOfStudent[s[index].StudentID])
	}
}

// UserIDs returns the list of StudentGuardians user ids.
func (s StudentGuardians) UserIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(s))

	for _, studentGuardian := range s {
		list = append(list, studentGuardian.UserID)
	}

	return list
}

// SetUser sets User info in StudentGuardian model.
func (s *StudentGuardian) SetUser(user User) {
	s.User = user
}

// SetUsers sets user info into list of StudentGuardians.
func (s StudentGuardians) SetUsers(users Users) {
	mapOfUsers := make(map[uuid.UUID]User, len(users))
	for index := range users {
		mapOfUsers[users[index].ID] = users[index]
	}

	for index := range s {
		s[index].SetUser(mapOfUsers[s[index].UserID])
	}
}

// StudentGuardianRelation represents the student guardian's relation.
type StudentGuardianRelation string

const (
	// StudentGuardianRelationTypeMother student guardian relation type mother.
	StudentGuardianRelationTypeMother StudentGuardianRelation = "mother"
	// StudentGuardianRelationTypeFather student guardian relation type father.
	StudentGuardianRelationTypeFather StudentGuardianRelation = "father"
	// StudentGuardianRelationTypeGuardian student guardian relation type guardian.
	StudentGuardianRelationTypeGuardian StudentGuardianRelation = "guardian"
	// StudentGuardianRelationTypeRelative student guardian relation type relative.
	StudentGuardianRelationTypeRelative StudentGuardianRelation = "relative"
)

// Validate validate student guardian relation.
func (s StudentGuardianRelation) Validate() bool {
	switch s {
	case StudentGuardianRelationTypeMother, StudentGuardianRelationTypeFather,
		StudentGuardianRelationTypeGuardian, StudentGuardianRelationTypeRelative:
		return true
	}

	return false
}

// StudentGuardianListFilter filter for the list of guardians.
type StudentGuardianListFilter struct {
	ListFilter

	CreatedDate     DateFilter
	GroupIDs        []uuid.UUID
	SchoolIDs       []uuid.UUID
	OrganizationIDs []uuid.UUID
}

// NewStudentGuardianListFilter creates a new StudentGuardianListFilter domain.
func NewStudentGuardianListFilter(
	list ListFilter,
	createdDate DateFilter,
	groupIDs []uuid.UUID,
	schoolIDs []uuid.UUID,
	organizationIDs []uuid.UUID,
) StudentGuardianListFilter {
	return StudentGuardianListFilter{
		ListFilter:      list,
		CreatedDate:     createdDate,
		GroupIDs:        groupIDs,
		SchoolIDs:       schoolIDs,
		OrganizationIDs: organizationIDs,
	}
}
