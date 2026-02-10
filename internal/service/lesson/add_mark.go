package lesson

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// AddMarkArgs is mark arguments for adding.
type AddMarkArgs struct {
	LessonID    uuid.UUID
	StudentID   uuid.UUID
	Mark        string
	Description *string
}

// AddMark adds a new mark to student.
func (s *Service) AddMark(ctx context.Context, args AddMarkArgs) (domain.Mark, error) {
	markDomain := domain.NewMark(args.LessonID, args.StudentID, args.Mark, args.Description, s.now)

	err := s.lessonRepo.AddMark(ctx, markDomain)
	if err != nil {
		return domain.Mark{}, fmt.Errorf("failed to add mark: %w", err)
	}

	markDomain, err = s.MarkByID(ctx, markDomain.ID)
	if err != nil {
		return domain.Mark{}, fmt.Errorf("failed to get mark by id: %w", err)
	}

	return markDomain, nil
}
