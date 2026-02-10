package populate

import (
	"bum-service/internal/infrastructure/repository"
	"bum-service/pkg/liblog"
	"bum-service/pkg/postgres"
	"bum-service/pkg/transaction"
	"context"
	"fmt"
)

func (s *Service) databaseService(ctx context.Context) (service *postgres.SqlxDBClient, err error) {
	var (
		logger = liblog.Must(ctx)
		url    = s.cfg.Infrastructure.Database.GetDatabaseDSN()
	)

	logger.Info("Postgres connection to database")

	dbMaxConnections := 100
	pg, err := postgres.NewClient(ctx, url, postgres.WithMaxOpenConnections(dbMaxConnections))
	if err != nil {
		logger.Errorf("Postgres connection %v", err)
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// adding http service in order to close it gracefully later
	s.services.Add(pg)

	logger.Info("database connection was established")

	s.container.Repo.db = pg

	return pg, nil
}

func (s *Service) db() postgres.DB {
	return s.container.Repo.db
}

func (s *Service) sessionAdapter() *transaction.SessionAdapter {
	if s.container.Repo.sessionAdapter != nil {
		return s.container.Repo.sessionAdapter
	}

	s.container.Repo.sessionAdapter = transaction.NewSessionAdapter(s.db())

	return s.container.Repo.sessionAdapter
}

func (s *Service) systemRepository() *repository.System {
	if s.container.Repo.systemRepository != nil {
		return s.container.Repo.systemRepository
	}

	s.container.Repo.systemRepository = repository.NewSystem(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.systemRepository
}

func (s *Service) gradesRepository() *repository.Grades {
	if s.container.Repo.gradesRepository != nil {
		return s.container.Repo.gradesRepository
	}

	s.container.Repo.gradesRepository = repository.NewGrades(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.gradesRepository
}

func (s *Service) eduOrganizationRepository() *repository.EduOrganization {
	if s.container.Repo.eduOrganizationRepository != nil {
		return s.container.Repo.eduOrganizationRepository
	}

	s.container.Repo.eduOrganizationRepository = repository.NewEduOrganization(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.eduOrganizationRepository
}

func (s *Service) ownerRepository() *repository.Owner {
	if s.container.Repo.ownerRepository != nil {
		return s.container.Repo.ownerRepository
	}

	s.container.Repo.ownerRepository = repository.NewOwner(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.ownerRepository
}

func (s *Service) schoolRepository() *repository.School {
	if s.container.Repo.schoolRepository != nil {
		return s.container.Repo.schoolRepository
	}

	s.container.Repo.schoolRepository = repository.NewSchool(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.schoolRepository
}

func (s *Service) groupRepository() *repository.Group {
	if s.container.Repo.groupRepository != nil {
		return s.container.Repo.groupRepository
	}

	s.container.Repo.groupRepository = repository.NewGroup(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.groupRepository
}

func (s *Service) groupSubjectsRepository() *repository.GroupSubjects {
	if s.container.Repo.groupSubjectsRepository != nil {
		return s.container.Repo.groupSubjectsRepository
	}

	s.container.Repo.groupSubjectsRepository = repository.NewGroupSubjects(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.groupSubjectsRepository
}

func (s *Service) userRepository() *repository.User {
	if s.container.Repo.userRepository != nil {
		return s.container.Repo.userRepository
	}

	s.container.Repo.userRepository = repository.NewUser(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.userRepository
}

func (s *Service) directorRepository() *repository.Director {
	if s.container.Repo.directorRepository != nil {
		return s.container.Repo.directorRepository
	}

	s.container.Repo.directorRepository = repository.NewDirector(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.directorRepository
}

func (s *Service) headmasterRepository() *repository.Headmaster {
	if s.container.Repo.headmasterRepository != nil {
		return s.container.Repo.headmasterRepository
	}

	s.container.Repo.headmasterRepository = repository.NewHeadmaster(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.headmasterRepository
}

func (s *Service) subjectRepository() *repository.Subject {
	if s.container.Repo.subjectRepository != nil {
		return s.container.Repo.subjectRepository
	}

	s.container.Repo.subjectRepository = repository.NewSubject(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.subjectRepository
}

func (s *Service) teacherRepository() *repository.Teacher {
	if s.container.Repo.teacherRepository != nil {
		return s.container.Repo.teacherRepository
	}

	s.container.Repo.teacherRepository = repository.NewTeacher(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.teacherRepository
}

func (s *Service) lessonRepository() *repository.Lesson {
	if s.container.Repo.lessonRepository != nil {
		return s.container.Repo.lessonRepository
	}

	s.container.Repo.lessonRepository = repository.NewLesson(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.lessonRepository
}

func (s *Service) studentRepository() *repository.Student {
	if s.container.Repo.studentRepository != nil {
		return s.container.Repo.studentRepository
	}

	s.container.Repo.studentRepository = repository.NewStudent(
		s.db(),
		s.sessionAdapter(),
	)

	return s.container.Repo.studentRepository
}
