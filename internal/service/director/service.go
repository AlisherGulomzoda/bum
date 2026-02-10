package director

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is a director use case.
type Service struct {
	userService     IUserService
	userInfoService IUserInfoService
	schoolService   ISchoolService

	directorRepo IDirectorRepo

	sessionAdapter transaction.Session
	logger         liblog.Logger
	now            func() time.Time
}

// NewService creates a new director use case.
func NewService(
	userService IUserService,
	userInfoService IUserInfoService,
	schoolService ISchoolService,

	directorRepo IDirectorRepo,

	sessionAdapter transaction.Session,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		userService:     userService,
		userInfoService: userInfoService,
		schoolService:   schoolService,

		directorRepo: directorRepo,

		sessionAdapter: sessionAdapter,
		logger:         logger,
		now:            nowFunc,
	}
}
