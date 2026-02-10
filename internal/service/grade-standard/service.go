package grades

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is a grade-standard use case.
type Service struct {
	gradesRepo IGradesRepo

	sessionAdapter transaction.Session
	logger         liblog.Logger
	now            func() time.Time
}

// NewService creates a new grade-standard use case.
func NewService(
	gradesRepo IGradesRepo,

	sessionAdapter *transaction.SessionAdapter,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		gradesRepo: gradesRepo,

		sessionAdapter: sessionAdapter,
		logger:         logger,
		now:            nowFunc,
	}
}
