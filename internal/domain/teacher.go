package domain

import (
	"time"

	"github.com/google/uuid"
)

// Teacher is the structure of teacher.
type Teacher struct {
	ID              uuid.UUID
	RoleID          uuid.UUID
	UserID          uuid.UUID
	SchoolID        uuid.UUID
	SchoolShortInfo *SchoolShortInfo
	Phone           *string
	Email           *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	User
}

// NewTeacher creates a new Teacher domain.
func NewTeacher(
	roleID uuid.UUID,
	userID uuid.UUID,
	schoolID uuid.UUID,
	phone *string,
	email *string,
	nowFunc func() time.Time,
) Teacher {
	now := nowFunc()

	return Teacher{
		ID:       uuid.New(),
		RoleID:   roleID,
		UserID:   userID,
		SchoolID: schoolID,
		Phone:    phone,
		Email:    email,
		
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetUser sets User info in Teacher model.
func (t *Teacher) SetUser(user User) {
	t.User = user
}

// SetShortSchool sets school info into director.
func (t *Teacher) SetShortSchool(schoolShort SchoolShortInfo) {
	t.SchoolShortInfo = &schoolShort
}

// Teachers are collection of Teacher.
type Teachers []Teacher

// UserIDs returns the list of Teacher users ids.
func (t Teachers) UserIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(t))

	for _, teacher := range t {
		list = append(list, teacher.UserID)
	}

	return list
}

// SchoolIDs returns the list of Teachers schools ids.
func (t Teachers) SchoolIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(t))

	for _, teacher := range t {
		list = append(list, teacher.SchoolID)
	}

	return list
}

// SetUsers sets users info into list of teachers.
func (t Teachers) SetUsers(users Users) {
	mapOfUsers := make(map[uuid.UUID]User, len(users))
	for _, user := range users {
		mapOfUsers[user.ID] = user
	}

	for index := range t {
		t[index].SetUser(mapOfUsers[t[index].UserID])
	}
}

// SetShortSchools sets school info into list of directors.
func (t Teachers) SetShortSchools(schoolShorts SchoolShortInfos) {
	mapOfSchoolShort := make(map[uuid.UUID]SchoolShortInfo, len(schoolShorts))
	for _, schoolShort := range schoolShorts {
		mapOfSchoolShort[schoolShort.ID] = schoolShort
	}

	for index := range t {
		t[index].SetShortSchool(mapOfSchoolShort[t[index].SchoolID])
	}
}

// TeacherListFilter filter for the list of teachers.
type TeacherListFilter struct {
	ListFilter

	CreatedDate     DateFilter
	GroupIDs        []uuid.UUID
	SchoolIDs       []uuid.UUID
	OrganizationIDs []uuid.UUID
}

// NewTeacherListFilter creates a new TeacherListFilter domain.
func NewTeacherListFilter(
	listFilter ListFilter,
	createdDate DateFilter,
	groupIDs []uuid.UUID,
	schoolIDs []uuid.UUID,
	organizationIDs []uuid.UUID,
) TeacherListFilter {
	return TeacherListFilter{
		ListFilter:      listFilter,
		CreatedDate:     createdDate,
		GroupIDs:        groupIDs,
		SchoolIDs:       schoolIDs,
		OrganizationIDs: organizationIDs,
	}
}
