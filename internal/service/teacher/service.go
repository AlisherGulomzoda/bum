package teacher

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is teacher use case.
type Service struct {
	userService     IUserService
	userInfoService IUserInfoService
	schoolService   ISchoolService

	teacherRepo ITeacherRepo

	sessionAdapter transaction.Session
	logger         liblog.Logger
	now            func() time.Time
}

// NewService creates a new teacher use case.
func NewService(
	userService IUserService,
	userInfoService IUserInfoService,
	schoolService ISchoolService,

	teacherRepo ITeacherRepo,

	sessionAdapter transaction.Session,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		userService:     userService,
		userInfoService: userInfoService,
		schoolService:   schoolService,

		teacherRepo: teacherRepo,

		sessionAdapter: sessionAdapter,
		logger:         logger,
		now:            nowFunc,
	}
}
