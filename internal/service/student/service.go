package student

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is student service.
type Service struct {
	userService     IUserService
	groupService    IGroupService
	userInfoService IUserInfoService
	schoolService   ISchoolService

	studentRepo IStudentRepo

	sessionAdapter transaction.Session
	logger         liblog.Logger
	now            func() time.Time
}

// NewService creates a new Student service.
func NewService(
	userService IUserService,
	groupService IGroupService,
	userInfoService IUserInfoService,
	schoolService ISchoolService,

	studentRepo IStudentRepo,

	sessionAdapter transaction.Session,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		userService:     userService,
		groupService:    groupService,
		userInfoService: userInfoService,
		schoolService:   schoolService,

		studentRepo: studentRepo,

		sessionAdapter: sessionAdapter,
		logger:         logger,
		now:            nowFunc,
	}
}
