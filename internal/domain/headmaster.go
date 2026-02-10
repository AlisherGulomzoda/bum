package domain

import (
	"time"

	"github.com/google/uuid"
)

// Headmaster is structure of Headmaster.
type Headmaster struct {
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

// NewHeadmaster creates a new Headmaster domain.
func NewHeadmaster(
	roleID uuid.UUID,
	userID uuid.UUID,
	schoolID uuid.UUID,
	phone *string,
	email *string,
	nowFunc func() time.Time,
) Headmaster {
	now := nowFunc()

	return Headmaster{
		ID:        uuid.New(),
		RoleID:    roleID,
		UserID:    userID,
		SchoolID:  schoolID,
		Phone:     phone,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetUser sets User info in Headmaster model.
func (h *Headmaster) SetUser(user User) {
	h.User = user
}

// SetShortSchool sets school info into headmaster.
func (h *Headmaster) SetShortSchool(schoolShort SchoolShortInfo) {
	h.SchoolShortInfo = &schoolShort
}

// HeadmasterListFilter filter for the list of Headmaster.
type HeadmasterListFilter struct {
	CreatedDate DateFilter
	ListFilter
	SchoolIDs []uuid.UUID
}

// NewHeadmasterListFilter creates a new HeadmasterListFilter domain.
func NewHeadmasterListFilter(createdDate DateFilter, list ListFilter, schoolIDs []uuid.UUID) HeadmasterListFilter {
	return HeadmasterListFilter{
		CreatedDate: createdDate,
		ListFilter:  list,
		SchoolIDs:   schoolIDs,
	}
}

// Headmasters is a list of Headmaster.
type Headmasters []Headmaster

// UserIDs returns the list of Director users ids.
func (h Headmasters) UserIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(h))

	for _, headmaster := range h {
		list = append(list, headmaster.UserID)
	}

	return list
}

// SetUsers sets users info into list of directors.
func (h Headmasters) SetUsers(users Users) {
	mapOfUsers := make(map[uuid.UUID]User, len(users))
	for _, user := range users {
		mapOfUsers[user.ID] = user
	}

	for index := range h {
		h[index].SetUser(mapOfUsers[h[index].UserID])
	}
}

// SchoolIDs returns the list of Headmaster schools ids.
func (h Headmasters) SchoolIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(h))

	for _, headmaster := range h {
		list = append(list, headmaster.SchoolID)
	}

	return list
}

// SetShortSchools sets school info into list of headmasters.
func (h Headmasters) SetShortSchools(schoolShorts SchoolShortInfos) {
	mapOfSchoolShort := make(map[uuid.UUID]SchoolShortInfo, len(schoolShorts))
	for _, schoolShort := range schoolShorts {
		mapOfSchoolShort[schoolShort.ID] = schoolShort
	}

	for index := range h {
		h[index].SetShortSchool(mapOfSchoolShort[h[index].SchoolID])
	}
}
