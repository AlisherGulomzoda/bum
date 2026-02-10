package system

import (
	"context"
	"fmt"
)

// HealthCheck checks whether the system is healthy.
func (s Service) HealthCheck(ctx context.Context) error {
	if err := s.systemRepo.Ping(ctx); err != nil {
		return fmt.Errorf("failed to check system health: %w", err)
	}

	return nil
}
