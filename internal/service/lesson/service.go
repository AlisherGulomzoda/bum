package lesson

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is a headmaster use case.
type Service struct {
	schoolService ISchoolService
	groupService  IGroupService

	lessonRepo ILessonRepo

	sessionAdapter transaction.Session
	logger         liblog.Logger
	now            func() time.Time
}

// NewService creates a new headmaster use case.
func NewService(
	schoolService ISchoolService,
	groupService IGroupService,

	lessonRepo ILessonRepo,

	sessionAdapter *transaction.SessionAdapter,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		schoolService: schoolService,
		groupService:  groupService,

		lessonRepo: lessonRepo,

		sessionAdapter: sessionAdapter,
		logger:         logger,
		now:            nowFunc,
	}
}
