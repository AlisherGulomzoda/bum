package populate

import (
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
)

func (s *Service) authService() *auth.Service {
	if s.container.Service.authService.Service != nil {
		return s.container.Service.authService.Service
	}

	s.container.Service.authService.Service = auth.NewService(
		s.container.Service.userService,

		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.authService.Service
}

func (s *Service) systemService() *system.Service {
	if s.container.Service.systemService.Service != nil {
		return s.container.Service.systemService.Service
	}

	s.container.Service.systemService.Service = system.NewService(
		s.systemRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.systemService.Service
}

func (s *Service) eduOrganizationService() *eduorganization.Service {
	if s.container.Service.eduOrganizationService.Service != nil {
		return s.container.Service.eduOrganizationService.Service
	}

	s.container.Service.eduOrganizationService.Service = eduorganization.NewService(
		s.eduOrganizationRepository(),

		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.eduOrganizationService.Service
}

func (s *Service) ownerService() *owner.Service {
	if s.container.Service.ownerService.Service != nil {
		return s.container.Service.ownerService.Service
	}

	s.container.Service.ownerService.Service = owner.NewService(
		s.container.Service.userService,
		s.container.Service.userInfoService,

		s.ownerRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.ownerService.Service
}

func (s *Service) schoolService() *school.Service {
	if s.container.Service.schoolService.Service != nil {
		return s.container.Service.schoolService.Service
	}

	s.container.Service.schoolService.Service = school.NewService(
		s.container.Service.eduOrganizationService,
		s.container.Service.gradesService,
		s.container.Service.teacherService,
		s.container.Service.studentService,

		s.schoolRepository(),
		s.groupRepository(),
		s.groupSubjectsRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.schoolService.Service
}

func (s *Service) userService() *user.Service {
	if s.container.Service.userService.Service != nil {
		return s.container.Service.userService.Service
	}

	s.container.Service.userService.Service = user.NewService(
		s.container.Service.eduOrganizationService,
		s.container.Service.schoolService,
		s.container.Service.groupService,

		s.userRepository(),
		s.studentRepository(),
		s.cfg.Application.PasswordCost,

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.userService.Service
}

func (s *Service) userInfoService() *userinfo.Service {
	if s.container.Service.userInfoService.Service != nil {
		return s.container.Service.userInfoService.Service
	}

	s.container.Service.userInfoService.Service = userinfo.NewService(
		s.userRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.userInfoService.Service
}

func (s *Service) directorService() *director.Service {
	if s.container.Service.directorService.Service != nil {
		return s.container.Service.directorService.Service
	}

	s.container.Service.directorService.Service = director.NewService(
		s.container.Service.userService,
		s.container.Service.userInfoService,
		s.container.Service.schoolService,

		s.directorRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.directorService.Service
}

func (s *Service) headmasterService() *headmaster.Service {
	if s.container.Service.headmasterService.Service != nil {
		return s.container.Service.headmasterService.Service
	}

	s.container.Service.headmasterService.Service = headmaster.NewService(
		s.container.Service.userService,
		s.container.Service.userInfoService,
		s.container.Service.schoolService,

		s.headmasterRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.headmasterService.Service
}

func (s *Service) subjectService() *subject.Service {
	if s.container.Service.subjectService.Service != nil {
		return s.container.Service.subjectService.Service
	}

	s.container.Service.subjectService.Service = subject.NewService(
		s.subjectRepository(),

		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.subjectService.Service
}

func (s *Service) teacherService() *teacher.Service {
	if s.container.Service.teacherService.Service != nil {
		return s.container.Service.teacherService.Service
	}

	s.container.Service.teacherService.Service = teacher.NewService(
		s.container.Service.userService,
		s.container.Service.userInfoService,
		s.container.Service.schoolService,

		s.teacherRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.teacherService.Service
}

func (s *Service) studentService() *student.Service {
	if s.container.Service.studentService.Service != nil {
		return s.container.Service.studentService.Service
	}

	s.container.Service.studentService.Service = student.NewService(
		s.container.Service.userService,
		s.container.Service.groupService,
		s.container.Service.userInfoService,
		s.container.Service.schoolService,

		s.studentRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.studentService.Service
}

func (s *Service) gradesService() *grades.Service {
	if s.container.Service.gradesService.Service != nil {
		return s.container.Service.gradesService.Service
	}

	s.container.Service.gradesService.Service = grades.NewService(
		s.gradesRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.gradesService.Service
}

func (s *Service) groupService() *school.Service {
	if s.container.Service.groupService.Service != nil {
		return s.container.Service.groupService.Service
	}

	s.container.Service.groupService.Service = s.schoolService()

	return s.container.Service.groupService.Service
}

func (s *Service) lessonService() *lesson.Service {
	if s.container.Service.lessonService.Service != nil {
		return s.container.Service.lessonService.Service
	}

	s.container.Service.lessonService.Service = lesson.NewService(
		s.container.Service.groupService,
		s.container.Service.schoolService,

		s.lessonRepository(),

		s.sessionAdapter(),
		s.logger(),
		s.nowFunc(),
	)

	return s.container.Service.lessonService.Service
}
