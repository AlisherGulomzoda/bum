package headmaster

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is a headmaster use case.
type Service struct {
	userService     IUserService
	userInfoService IUserInfoService
	schoolService   ISchoolService

	headmasterRepo IHeadmasterRepo

	sessionAdapter transaction.Session
	logger         liblog.Logger
	now            func() time.Time
}

// NewService creates a new headmaster use case.
func NewService(
	userService IUserService,
	userInfoService IUserInfoService,
	schoolService ISchoolService,

	headmasterRepo IHeadmasterRepo,

	sessionAdapter *transaction.SessionAdapter,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		userService:     userService,
		userInfoService: userInfoService,
		schoolService:   schoolService,

		headmasterRepo: headmasterRepo,

		sessionAdapter: sessionAdapter,
		logger:         logger,
		now:            nowFunc,
	}
}
