package headmaster

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// HeadmasterByID get headmaster by id.
func (s Service) HeadmasterByID(ctx context.Context, id uuid.UUID) (domain.Headmaster, error) {
	headmaster, err := s.headmasterRepo.HeadmasterByIDTx(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to get headmaster by id: %w", err)
		return headmaster, err
	}

	user, err := s.userInfoService.UserByID(ctx, headmaster.UserID)
	if err != nil {
		err = fmt.Errorf("failed to get user by id: %w", err)
		return headmaster, err
	}

	headmaster.SetUser(user)

	return headmaster, nil
}
