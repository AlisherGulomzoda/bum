package domain

import (
	"time"

	"github.com/google/uuid"
)

// Director is the structure of Director.
type Director struct {
	ID              uuid.UUID
	RoleID          uuid.UUID
	UserID          uuid.UUID
	SchoolID        uuid.UUID
	SchoolShortInfo *SchoolShortInfo
	Phone           *string
	Email           *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
	User
}

// NewDirector creates a new Director domain.
func NewDirector(
	roleID uuid.UUID,
	userID uuid.UUID,
	schoolID uuid.UUID,
	phone *string,
	email *string,
	nowFunc func() time.Time,
) Director {
	now := nowFunc()

	return Director{
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

// SetUser sets User info in Director model.
func (d *Director) SetUser(user User) {
	d.User = user
}

// SetShortSchool sets school info into director.
func (d *Director) SetShortSchool(schoolShort SchoolShortInfo) {
	d.SchoolShortInfo = &schoolShort
}

// DirectorListFilter filter for the list of directors.
type DirectorListFilter struct {
	ListFilter
	CreatedDate DateFilter
	SchoolIDs   []uuid.UUID
}

// NewDirectorListFilter creates a new DirectorListFilter domain.
func NewDirectorListFilter(createdDate DateFilter, list ListFilter, schoolIDs []uuid.UUID) DirectorListFilter {
	return DirectorListFilter{
		CreatedDate: createdDate,
		ListFilter:  list,
		SchoolIDs:   schoolIDs,
	}
}

// Directors are collection of Director.
type Directors []Director

// UserIDs returns the list of Directors users ids.
func (d Directors) UserIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(d))

	for _, director := range d {
		list = append(list, director.UserID)
	}

	return list
}

// SchoolIDs returns the list of Directors schools ids.
func (d Directors) SchoolIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(d))

	for _, director := range d {
		list = append(list, director.SchoolID)
	}

	return list
}

// SetUsers sets users info into list of directors.
func (d Directors) SetUsers(users Users) {
	mapOfUsers := make(map[uuid.UUID]User, len(users))
	for _, user := range users {
		mapOfUsers[user.ID] = user
	}

	for index := range d {
		d[index].SetUser(mapOfUsers[d[index].UserID])
	}
}

// SetShortSchools sets school info into list of directors.
func (d Directors) SetShortSchools(schoolShorts SchoolShortInfos) {
	mapOfSchoolShort := make(map[uuid.UUID]SchoolShortInfo, len(schoolShorts))
	for _, schoolShort := range schoolShorts {
		mapOfSchoolShort[schoolShort.ID] = schoolShort
	}

	for index := range d {
		d[index].SetShortSchool(mapOfSchoolShort[d[index].SchoolID])
	}
}
