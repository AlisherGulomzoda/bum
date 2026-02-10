package owner

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is owner use case.
type Service struct {
	userService     IUserService
	userInfoService IUserInfoService

	ownerRepo IOwnerRepository

	sessionAdapter transaction.Session
	logger         liblog.Logger
	now            func() time.Time
}

// NewService creates a new owner use case.
func NewService(
	userService IUserService,
	userInfoService IUserInfoService,
	ownerRepo IOwnerRepository,

	sessionAdapter *transaction.SessionAdapter,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		userService:     userService,
		userInfoService: userInfoService,
		ownerRepo:       ownerRepo,

		sessionAdapter: sessionAdapter,
		logger:         logger,
		now:            nowFunc,
	}
}
