package lesson

import (
	"context"
	"time"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// IUserInfoService represents a user info service for headmaster use cases.
type IUserInfoService interface {
	UserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UsersByIDs(ctx context.Context, ids []uuid.UUID) (domain.Users, error)
}

// ISchoolService represents a school service.
type ISchoolService interface {
	SchoolShortByIDs(ctx context.Context, ids []uuid.UUID) (domain.SchoolShortInfos, error)
}

// ILessonRepo is lesson repository.
type ILessonRepo interface {
	AssignsLessons(
		ctx context.Context, groupID uuid.UUID, firstDayOfWeek, firstDayOfNextWeek time.Time, lessons domain.Lessons,
	) error
	LessonsListTx(ctx context.Context, filters domain.LessonsListFilter) (domain.Lessons, error)

	AddMark(ctx context.Context, m domain.Mark) error
	MarkByIDTx(ctx context.Context, id uuid.UUID) (domain.Mark, error)
}

// IGroupService is a group service use case interface.
type IGroupService interface {
	GroupByID(ctx context.Context, groupID uuid.UUID) (domain.Group, error)
	GroupSubjectList(ctx context.Context, groupID uuid.UUID) (domain.GroupSubjects, error)
}
