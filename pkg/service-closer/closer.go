package closer

import (
	"context"

	"bum-service/pkg/liblog"
)

// Closer represents service interface.
type Closer interface {
	Close(ctx context.Context) error
}

// ServiceCloser represents service close interface.
type ServiceCloser interface {
	Close(ctx context.Context)
	Add(c Closer)
}

// Client is a client for list of services.
type Client struct {
	services []Closer
}

// NewCloser creates a new Closer client.
func NewCloser() Client {
	return Client{}
}

// Add adds a new service to list.
func (s *Client) Add(c Closer) {
	s.services = append(s.services, c)
}

// Close closes all services.
func (s *Client) Close(ctx context.Context) {
	logger := liblog.Must(ctx)

	for serviceOrder := len(s.services) - 1; serviceOrder >= 0; serviceOrder-- {
		if err := s.services[serviceOrder].Close(ctx); err != nil {
			logger.Errorf("failed to close service %+v", err)
		}
	}
}
