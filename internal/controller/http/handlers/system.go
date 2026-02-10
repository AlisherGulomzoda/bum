package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// System is system handlers.
type System struct {
	systemService ISystemService
}

// NewSystem creates a new system handler.
func NewSystem(service ISystemService) *System {
	return &System{
		systemService: service,
	}
}

// HealthCheck checks whether the system is healthy.
func (s System) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	if err := s.systemService.HealthCheck(ctx); err != nil {
		c.Status(http.StatusServiceUnavailable)

		return
	}

	c.Status(http.StatusOK)
}
