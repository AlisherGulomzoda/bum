package school

import (
	"time"

	"bum-service/pkg/liblog"
	"bum-service/pkg/transaction"
)

// Service is a school use case.
type Service struct {
	eduOrganizationService IEduOrganizationService
	gradeService           IGradeService
	teacherService         ITeacherService
	studentService         IStudentService

	schoolRepo        ISchoolRepo
	groupRepo         IGroupRepo
	groupSubjectsRepo IGroupSubjectsRepo

	sessionAdapter transaction.Session
	logger         liblog.Logger
	now            func() time.Time
}

// NewService creates a new school use case.
func NewService(
	eduOrganizationService IEduOrganizationService,
	gradeService IGradeService,
	teacherService ITeacherService,
	studentService IStudentService,

	schoolRepo ISchoolRepo,
	groupRepo IGroupRepo,
	groupSubjectsRepo IGroupSubjectsRepo,

	sessionAdapter transaction.Session,
	logger liblog.Logger,
	nowFunc func() time.Time,
) *Service {
	return &Service{
		eduOrganizationService: eduOrganizationService,
		gradeService:           gradeService,
		teacherService:         teacherService,
		studentService:         studentService,

		schoolRepo:        schoolRepo,
		groupRepo:         groupRepo,
		groupSubjectsRepo: groupSubjectsRepo,

		sessionAdapter: sessionAdapter,
		logger:         logger,
		now:            nowFunc,
	}
}
