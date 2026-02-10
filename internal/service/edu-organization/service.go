package eduorganization

import (
	"time"

	"bum-service/pkg/liblog"
)

// Service is an educational organization use case.
type Service struct {
	eduOrganizationRepo IEduOrganizationRepo

	logger liblog.Logger
	now    func() time.Time
}

// NewService creates a new educational organization use case.
func NewService(
	eduOrgRepo IEduOrganizationRepo,

	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		eduOrganizationRepo: eduOrgRepo,

		logger: logger,
		now:    nowFunc,
	}
}
