package lesson

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/utils"
)

// AddWeekLessonsArgs is args for adding lesson for a week.
type AddWeekLessonsArgs struct {
	SchoolID uuid.UUID
	GroupID  uuid.UUID
	WeekDate time.Time
	Lessons  []Lesson
}

// Lesson is lesson item for adding.
type Lesson struct {
	GroupSubjectID uuid.UUID
	TeacherID      *uuid.UUID
	AuditoriumID   uuid.UUID
	StartTime      time.Time
	EndTime        time.Time
	Description    *string
}

// AssignLessons assigns lessons to group for a weak.
func (s *Service) AssignLessons(ctx context.Context, args AddWeekLessonsArgs) (domain.Lessons, error) {
	// TODO: учитывать часовой пояс в будущем.
	var (
		firstDayOfWeek     = utils.FirstDayOfWeek(args.WeekDate)
		firstDayOfNextWeek = firstDayOfWeek.AddDate(0, 0, utils.WeekDaysCount)
		lessons            = make(domain.Lessons, 0, len(args.Lessons))
	)

	group, err := s.groupService.GroupByID(ctx, args.GroupID)
	if err != nil {
		return domain.Lessons{}, fmt.Errorf("failed to get group by id and school id: %w", err)
	}

	groupSubject, err := s.groupService.GroupSubjectList(ctx, args.GroupID)
	if err != nil {
		return domain.Lessons{}, fmt.Errorf("failed to get group subject list: %w", err)
	}

	groupSubjectMap := groupSubject.MapByID()

	for _, l := range args.Lessons {
		lessonDomain := domain.NewLesson(
			group.SchoolID,
			l.GroupSubjectID,
			l.TeacherID,
			groupSubjectMap[l.GroupSubjectID].TeacherID,
			l.AuditoriumID,
			l.StartTime,
			l.EndTime,
			l.Description,

			s.now,
		)

		lessons = append(lessons, lessonDomain)
	}

	err = s.lessonRepo.AssignsLessons(ctx, args.GroupID, firstDayOfWeek, firstDayOfNextWeek, lessons)
	if err != nil {
		return nil, fmt.Errorf("failed to assign lessons: %w", err)
	}

	return lessons, nil
}
