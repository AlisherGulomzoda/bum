package owner

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// OwnerByID get owner by id.
func (s Service) OwnerByID(ctx context.Context, id uuid.UUID) (domain.Owner, error) {
	owner, err := s.ownerRepo.OwnerByIDTx(ctx, id)
	if err != nil {
		return owner, fmt.Errorf("failed to get owner by id: %w", err)
	}

	user, err := s.userInfoService.UserByID(ctx, owner.UserID)
	if err != nil {
		return owner, fmt.Errorf("failed to get user by id: %w", err)
	}

	owner.SetUser(user)

	return owner, nil
}

// OwnerByUserIDAndSchoolID get owner by user_id and school_id.
func (s Service) OwnerByUserIDAndSchoolID(ctx context.Context, schoolID, userID uuid.UUID) (domain.Owner, error) {
	owner, err := s.ownerRepo.OwnerByUserIDAndSchoolIDTx(ctx, schoolID, userID)
	if err != nil {
		return owner, fmt.Errorf("failed to get owner by user_id and school_id: %w", err)
	}

	user, err := s.userInfoService.UserByID(ctx, owner.UserID)
	if err != nil {
		return owner, fmt.Errorf("failed to get user by user_id and school_id: %w", err)
	}

	owner.SetUser(user)

	return owner, nil
}
