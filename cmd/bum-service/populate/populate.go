//nolint:gosec,gochecknoglobals //it's ok for generating test data.
package populate

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"

	"bum-service/config"
	"bum-service/internal/domain"
	"bum-service/pkg/liblog"
	closer "bum-service/pkg/service-closer"
)

const (
	env = "local"

	timeOutForGeneration = time.Minute * 10

	configPathFlag    = "config_path"
	defaultConfigPath = "config.yml"
)

var cfg *config.Config

// Cmd is a command for populating for test data.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "populate",
		Short:   "populate service with test data",
		PreRunE: loadConfigs,
		RunE:    runPopulate(),
	}

	cmd.Flags().String(configPathFlag, defaultConfigPath, "path to config yml file")
	return cmd
}

func loadConfigs(cmd *cobra.Command, _ []string) (err error) {
	configPath, err := cmd.Flags().GetString(configPathFlag)
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	cfg, err = config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return nil
}

// Service is a service structure.
type Service struct {
	cfg *config.Config

	container Container

	services     closer.Client
	gradesList   map[string]domain.GradeStandard
	subjectsList map[string]domain.Subject
	randomUsers  domain.Users
}

// Close closes the service.
func (s *Service) Close() {
	logger := s.logger()

	logger.Info("Closing the service")

	ctx := liblog.With(context.Background(), logger)

	s.services.Close(ctx)

	logger.Info("service closed")
}

// runPopulate runs the Server.
func runPopulate() func(_ *cobra.Command, _ []string) error {
	return func(_ *cobra.Command, _ []string) error {

		if cfg.Application.Env != env {
			return fmt.Errorf("allowed only for %s environment", env)
		}

		now := time.Now()
		defer fmt.Println("populate time: ", time.Since(now).Seconds())

		s := Service{
			cfg: cfg,

			services: closer.NewCloser(),

			container: Container{
				Service: NewServiceContainer(),
			},

			gradesList:   make(map[string]domain.GradeStandard),
			subjectsList: make(map[string]domain.Subject),
		}

		s.container.logger = liblog.NewDummyLogger()

		defer s.Close()

		rand.New(rand.NewSource(time.Now().UnixNano()))

		ctx := liblog.With(context.Background(), s.logger())

		ctx, cancel := context.WithTimeout(ctx, timeOutForGeneration)
		defer cancel()

		_, err := s.databaseService(ctx)
		if err != nil {
			return fmt.Errorf("failed to get database: %w", err)
		}

		_ = s.sessionAdapter()

		_ = s.nowFunc()

		_ = s.systemService()

		_ = s.gradesService()

		_ = s.eduOrganizationService()

		s.ownerService()

		_ = s.schoolService()

		_ = s.userService()

		_ = s.authService()

		_ = s.userInfoService()

		_ = s.directorService()

		_ = s.headmasterService()

		_ = s.subjectService()

		_ = s.teacherService()

		_ = s.groupService()

		_ = s.studentService()

		_ = s.lessonService()

		if err = s.container.Service.CheckInitialized(); err != nil {
			return err
		}

		if err = s.container.Repo.CheckInitialized(); err != nil {
			return err
		}

		_, userCount, err := s.userService().UserList(
			ctx,
			domain.UserListFilter{
				ListFilter: domain.ListFilter{
					Pagination: domain.Pagination{
						Limit:  domain.PaginationDefaultLimit,
						Offset: domain.PaginationDefaultPage,
					},
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to get list of users: %w", err)
		}

		if userCount > 0 {
			return nil
		}

		// add random users.
		err = s.genRandomUsers(ctx)
		if err != nil {
			return err
		}

		// add grade standards.
		err = s.addGradeStandards(ctx)
		if err != nil {
			return err
		}

		// add subjects.
		err = s.addSubjects(ctx)
		if err != nil {
			return err
		}

		// generate multi user
		multiUser, err := s.generateUserWithMultipleRoles(ctx)
		if err != nil {
			return err
		}

		_ = multiUser

		// add organizations
		err = s.generateOrganizations(ctx)
		if err != nil {
			return err
		}

		return nil
	}
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
