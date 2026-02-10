package lesson

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// MarkByID returns mark by id.
func (s *Service) MarkByID(ctx context.Context, markID uuid.UUID) (domain.Mark, error) {
	mark, err := s.lessonRepo.MarkByIDTx(ctx, markID)
	if err != nil {
		return domain.Mark{}, fmt.Errorf("failed get mark by id from database: %w", err)
	}

	return mark, nil
}
