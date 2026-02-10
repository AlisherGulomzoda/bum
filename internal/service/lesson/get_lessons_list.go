package lesson

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// LessonsList returns lessons list.
func (s *Service) LessonsList(ctx context.Context, filters domain.LessonsListFilter) (domain.Lessons, error) {
	list, err := s.lessonRepo.LessonsListTx(ctx, filters)
	if err != nil {
		return domain.Lessons{}, fmt.Errorf("failed get lessons list from database: %w", err)
	}

	return list, nil
}
