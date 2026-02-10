//nolint:forbidigo,revive,goerr113 // it's ok
package app

import (
	"fmt"
	"reflect"
	"time"

	"bum-service/internal/infrastructure/repository"
	"bum-service/internal/service/auth"
	"bum-service/internal/service/director"
	eduorganization "bum-service/internal/service/edu-organization"
	grades "bum-service/internal/service/grade-standard"
	"bum-service/internal/service/headmaster"
	"bum-service/internal/service/lesson"
	"bum-service/internal/service/owner"
	"bum-service/internal/service/school"
	"bum-service/internal/service/student"
	"bum-service/internal/service/subject"
	"bum-service/internal/service/system"
	"bum-service/internal/service/teacher"
	"bum-service/internal/service/user"
	"bum-service/internal/service/userinfo"
	"bum-service/pkg/liblog"
	"bum-service/pkg/postgres"
	"bum-service/pkg/transaction"
)

// Container is dependency container to solve loop dependencies.
type Container struct {
	Service ServiceContainer

	Repo RepoContainer

	logger  liblog.Logger
	nowFunc func() time.Time
}

// ServiceContainer is service container.
type ServiceContainer struct {
	authService            *struct{ *auth.Service }
	systemService          *struct{ *system.Service }
	eduOrganizationService *struct{ *eduorganization.Service }
	ownerService           *struct{ *owner.Service }
	schoolService          *struct{ *school.Service }
	directorService        *struct{ *director.Service }
	headmasterService      *struct{ *headmaster.Service }
	subjectService         *struct{ *subject.Service }
	userService            *struct{ *user.Service }
	userInfoService        *struct{ *userinfo.Service }
	teacherService         *struct{ *teacher.Service }
	gradesService          *struct{ *grades.Service }
	studentService         *struct{ *student.Service }
	lessonService          *struct{ *lesson.Service }
	groupService           *struct{ *school.Service }
}

// NewServiceContainer creates a new service container.
func NewServiceContainer() ServiceContainer {
	return ServiceContainer{
		authService:            &struct{ *auth.Service }{},
		systemService:          &struct{ *system.Service }{},
		eduOrganizationService: &struct{ *eduorganization.Service }{},
		ownerService:           &struct{ *owner.Service }{},
		schoolService:          &struct{ *school.Service }{},
		directorService:        &struct{ *director.Service }{},
		headmasterService:      &struct{ *headmaster.Service }{},
		subjectService:         &struct{ *subject.Service }{},
		userService:            &struct{ *user.Service }{},
		userInfoService:        &struct{ *userinfo.Service }{},
		teacherService:         &struct{ *teacher.Service }{},
		gradesService:          &struct{ *grades.Service }{},
		studentService:         &struct{ *student.Service }{},
		lessonService:          &struct{ *lesson.Service }{},
		groupService:           &struct{ *school.Service }{},
	}
}

// CheckInitialized проверяет, что все поля структуры ServiceContainer не nil.
func (sc *ServiceContainer) CheckInitialized() error {
	v := reflect.ValueOf(sc).Elem()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			return fmt.Errorf("поле %s не инициализировано", v.Type().Field(i).Name)
		}
	}

	fmt.Println("okkkk")

	return nil
}

// RepoContainer is repo container.
type RepoContainer struct {
	db             postgres.DB
	sessionAdapter *transaction.SessionAdapter

	systemRepository          *repository.System
	eduOrganizationRepository *repository.EduOrganization
	ownerRepository           *repository.Owner
	schoolRepository          *repository.School
	groupRepository           *repository.Group
	gradesRepository          *repository.Grades
	groupSubjectsRepository   *repository.GroupSubjects
	userRepository            *repository.User
	directorRepository        *repository.Director
	headmasterRepository      *repository.Headmaster //nolint:ireturn // it's ok
	subjectRepository         *repository.Subject
	teacherRepository         *repository.Teacher
	lessonRepository          *repository.Lesson
	studentRepository         *repository.Student
}

// CheckInitialized проверяет, что все поля структуры RepoContainer не nil.
func (sc *RepoContainer) CheckInitialized() error {
	v := reflect.ValueOf(sc).Elem()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			return fmt.Errorf("поле %s не инициализировано", v.Type().Field(i).Name)
		}
	}

	fmt.Println("okkkk2")

	return nil
}

//nolint:ireturn // it's ok
func (s *Service) logger() liblog.Logger {
	return s.container.logger
}

func (*Service) nowFunc() func() time.Time {
	return func() time.Time {
		return time.Now().UTC()
	}
}
