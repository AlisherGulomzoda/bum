package user

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is a user use case.
type Service struct {
	organizationService IEduOrganizationService
	schoolService       ISchoolService
	groupService        IGroupService

	userRepo    IUserRepo
	studentRepo IStudentRepo

	passwordCost   int
	sessionAdapter transaction.Session
	logger         liblog.Logger
	nowFunc        func() time.Time
}

// NewService creates a new user use case.
func NewService(
	organizationService IEduOrganizationService,
	schoolService ISchoolService,
	groupService IGroupService,

	userRepo IUserRepo,
	studentRepo IStudentRepo,

	passwordCost int,
	sessionAdapter transaction.Session,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		schoolService:       schoolService,
		organizationService: organizationService,
		groupService:        groupService,

		userRepo:    userRepo,
		studentRepo: studentRepo,

		passwordCost:   passwordCost,
		sessionAdapter: sessionAdapter,
		logger:         logger,
		nowFunc:        nowFunc,
	}
}
