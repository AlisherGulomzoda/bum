package subject

import (
	"time"

	"bum-service/pkg/liblog"
)

// Service is a subject use case.
type Service struct {
	subjectRepo ISubjectRepo

	logger liblog.Logger
	now    func() time.Time
}

// NewService creates a new subject use case.
func NewService(
	subjectRepo ISubjectRepo,

	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		subjectRepo: subjectRepo,

		logger: logger,
		now:    nowFunc,
	}
}
