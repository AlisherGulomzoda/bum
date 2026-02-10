package app

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"bum-service/config"
	"bum-service/pkg/liblog"
	closer "bum-service/pkg/service-closer"
)

const (
	dbMaxConnections = 30

	httpIdleTimeout       = 15 * time.Second
	httpReadHeaderTimeout = 5 * time.Second
)

// Service is a service structure.
type Service struct {
	cfg *config.Config

	container Container

	services closer.Client
}

// NewService creates a new Service.
func NewService(cfg *config.Config) (*Service, error) {
	s := &Service{
		cfg:      cfg,
		services: closer.NewCloser(),

		container: Container{
			Service: NewServiceContainer(),
		},
	}

	logger, err := liblog.NewZeroLog(convertConfigLogToLibLogConfig(s.cfg.Logger))
	if err != nil {
		return nil, fmt.Errorf("failed to create a new loger: %w", err)
	}

	s.container.logger = logger

	return s, nil
}

// Run runs the service.
func (s *Service) Run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = s.getPanicReason(r)
		}
	}()

	logger := s.logger()

	logger.Info("Service application starting ...")

	ctx := liblog.With(context.Background(), s.logger())

	_, err = s.databaseService(ctx)
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	s.db()

	s.sessionAdapter()

	s.ownerService()

	s.nowFunc()

	s.systemService()

	s.gradesService()

	s.eduOrganizationService()

	s.ownerService()

	s.schoolService()

	s.userService()

	s.authService()

	s.userInfoService()

	s.directorService()

	s.headmasterService()

	s.subjectService()

	s.teacherService()

	s.groupService()

	s.studentService()

	s.lessonService()

	if err = s.container.Service.CheckInitialized(); err != nil {
		logger.Error("Ошибка:", err)
		return err
	}

	if err = s.container.Repo.CheckInitialized(); err != nil {
		logger.Error("Ошибка:", err)
		return err
	}

	err = s.HTTPService()
	if err != nil {
		logger.Errorf("failed to run http server: %v", err)
		return err
	}

	return nil
}

// Close closes the service.
func (s *Service) Close() {
	logger := s.logger()

	logger.Info("Closing the service")

	ctx := liblog.With(context.Background(), logger)

	s.services.Close(ctx)

	logger.Info("service closed")
}

var errPanic = errors.New("panic occurred")

func (s *Service) getPanicReason(r any) error {
	var (
		err   = errPanic
		stack string
	)

	if r != nil {
		switch x := r.(type) {
		case string:
			err = fmt.Errorf("%w: %s", err, x)
		case error:
			err = x
			// stack = StackTrace(err)
		default:
			err = fmt.Errorf("%w: %+v", err, x)
		}

		if stack == "" {
			stack = string(debug.Stack())
		}

		s.logger().WithFields(map[string]any{
			"error": err,
			"stack": stack,
		}).Error("run service failed and panic was recovered")

		return err
	}

	return nil
}

func convertConfigLogToLibLogConfig(cfg config.Logger) liblog.Config {
	outputs := make([]liblog.Output, 0, len(cfg.Outputs))

	for _, output := range cfg.Outputs {
		outputs = append(outputs, liblog.Output(output))
	}

	return liblog.Config{
		Level:           liblog.Level(cfg.Level),
		Outputs:         outputs,
		Formatter:       liblog.Format(cfg.Formatter),
		TimeStampFormat: cfg.TimeStampFormat,
		Caller:          cfg.Caller,
		Sentry: liblog.SentryConfig{
			DSN:         cfg.Sentry.DSN,
			Environment: cfg.Sentry.Environment,
			Release:     cfg.Sentry.Release,
		},
	}
}
