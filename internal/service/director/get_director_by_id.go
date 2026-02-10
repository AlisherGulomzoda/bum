package director

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// DirectorByID get director by id.
func (s Service) DirectorByID(ctx context.Context, id uuid.UUID) (domain.Director, error) {
	director, err := s.directorRepo.DirectorByIDTx(ctx, id)
	if err != nil {
		return director, fmt.Errorf("failed to get director by id: %w", err)
	}

	user, err := s.userInfoService.UserByID(ctx, director.UserID)
	if err != nil {
		return director, fmt.Errorf("failed to get user by id: %w", err)
	}

	director.SetUser(user)

	return director, nil
}
