package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/gin-gonic/gin"

	controllerhttp "bum-service/internal/controller/http"
	"bum-service/pkg/liblog"
)

// HTTPService is http controller service.
func (s *Service) HTTPService() error {
	var (
		c       HTTPController
		logger  = s.logger()
		addr    = net.JoinHostPort(s.cfg.Controller.HTTP.Address, strconv.Itoa(s.cfg.Controller.HTTP.Port))
		handler = gin.New()
	)

	logger.Info("Registering handlers ...")

	err := controllerhttp.RegisterHandlers(
		handler,
		s.logger(),

		s.cfg.Application.JwtSecret,
		s.cfg.Application.AccessTokenExp,
		s.cfg.Application.RefreshTokenExp,

		s.authService(),
		s.systemService(),
		s.eduOrganizationService(),
		s.ownerService(),
		s.schoolService(),
		s.directorService(),
		s.headmasterService(),
		s.subjectService(),
		s.userService(),
		s.teacherService(),
		s.gradesService(),
		s.studentService(),
		s.lessonService(),
	)
	if err != nil {
		return fmt.Errorf("failed to create a new HTTP controller: %w", err)
	}

	logger.Info("Handlers were registered")

	c.server = &http.Server{
		IdleTimeout:       httpIdleTimeout,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		Addr:              addr,
		Handler:           handler,
	}

	go func() {
		logger.Infof("Starting to listen on %s", addr)

		if err := c.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("failed to listen http Server: %w", err))
		}
	}()

	// adding http service in order to close it gracefully later
	s.services.Add(c)

	return nil
}

// HTTPController is an HTTP controller.
type HTTPController struct {
	server *http.Server
}

// Close closes the http server if it is running.
func (s HTTPController) Close(ctx context.Context) error {
	logger := liblog.Must(ctx)

	err := s.server.Shutdown(ctx)
	if err != nil {
		logger.Errorf("Failed to shutdown http Server: %v", err)
		return fmt.Errorf("shutting down http server failed: %w", err)
	}

	logger.Info("http server closed successfully")

	return nil
}

func (s *Service) registerPprof() {
	var (
		c      PprofHandler
		logger = s.logger()
		addr   = net.JoinHostPort(s.cfg.Application.PprofHost, strconv.Itoa(s.cfg.Application.PprofPort))
		mux    = http.NewServeMux()
	)

	logger.Info("Registering handlers ...")

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	c.server = &http.Server{
		IdleTimeout:       httpIdleTimeout,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		Addr:              addr,
		Handler:           mux,
	}

	go func() {
		logger.Infof("Starting to listen on %s for pprof", addr)

		if err := c.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("failed to listen http Server for pprof: %w", err))
		}
	}()

	// adding http service to close it gracefully later
	s.services.Add(c)
}

type PprofHandler struct {
	HTTPController
}
