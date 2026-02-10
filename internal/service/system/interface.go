package system

import "context"

// ISystemRepo represents a repository for system use cases.
type ISystemRepo interface {
	Ping(ctx context.Context) error
}
