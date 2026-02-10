package userinfo

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is a user use case.
type Service struct {
	userInfoRepo IUserInfoRepo

	sessionAdapter transaction.Session
	logger         liblog.Logger
	nowFunc        func() time.Time
}

// NewService creates a new user use case.
func NewService(
	userInfoRepo IUserInfoRepo,

	sessionAdapter transaction.Session,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		userInfoRepo: userInfoRepo,

		sessionAdapter: sessionAdapter,
		logger:         logger,
		nowFunc:        nowFunc,
	}
}
