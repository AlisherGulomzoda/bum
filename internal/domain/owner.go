package domain

import (
	"time"

	"github.com/google/uuid"
)

// Owner is structure of Owner.
type Owner struct {
	ID             uuid.UUID
	RoleID         uuid.UUID
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Phone          *string
	Email          *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	User
}

// NewOwner creates a new Owner domain.
func NewOwner(
	roleID uuid.UUID,
	userID uuid.UUID,
	organizationID uuid.UUID,
	phone *string,
	email *string,

	nowFunc func() time.Time,
) Owner {
	now := nowFunc()

	return Owner{
		ID:             uuid.New(),
		RoleID:         roleID,
		UserID:         userID,
		OrganizationID: organizationID,
		Phone:          phone,
		Email:          email,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// SetUser  sets user to owner.
func (o *Owner) SetUser(user User) {
	o.User = user
}

// Owners are slice of Owner.
type Owners []Owner

// UserIDs returns the list of Directors users ids.
func (o Owners) UserIDs() []uuid.UUID {
	list := make([]uuid.UUID, 0, len(o))

	for _, owner := range o {
		list = append(list, owner.UserID)
	}

	return list
}

// OwnerListFilter is structure of Owner list filter.
type OwnerListFilter struct {
	ListFilter

	CreatedDate     DateFilter
	OrganizationIDs []uuid.UUID
}

// NewOwnerListFilter creates a new OwnerListFilter domain.
func NewOwnerListFilter(createdDate DateFilter, list ListFilter, organizationIDs []uuid.UUID) OwnerListFilter {
	return OwnerListFilter{
		CreatedDate:     createdDate,
		ListFilter:      list,
		OrganizationIDs: organizationIDs,
	}
}

// SetUsers sets users info into list of directors.
func (o Owners) SetUsers(users Users) {
	mapOfUsers := make(map[uuid.UUID]User, len(users))
	for _, user := range users {
		mapOfUsers[user.ID] = user
	}

	for index := range o {
		o[index].SetUser(mapOfUsers[o[index].UserID])
	}
}
