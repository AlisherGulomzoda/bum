package auth

import (
	"time"

	"bum-service/pkg/liblog"
)

// Service is an auth use case.
type Service struct {
	userService IUserService

	logger liblog.Logger
	now    func() time.Time
}

// NewService creates a new auth use case.
func NewService(
	userService IUserService,

	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		userService: userService,

		logger: logger,
		now:    nowFunc,
	}
}
