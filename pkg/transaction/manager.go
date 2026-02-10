//nolint:ireturn // we are writing wrapper over two different structures, and we cannot avoid using interface.
package transaction

import (
	"bum-service/pkg/postgres"
)

// Manager управляет транзакцией.
type Manager struct {
	rollback func() error
	commit   func() error
}

// Rollback откатывает транзакцию.
func (s *Manager) Rollback() error {
	return s.rollback()
}

// Commit завершает транзакцию.
func (s *Manager) Commit() error {
	return s.commit()
}

//nolint:wrapcheck // don't need to wrap
func newSolver(tx postgres.SqlxTxClient) *Manager {
	return &Manager{
		rollback: tx.Rollback,
		commit:   tx.Commit,
	}
}

func noopSolver() *Manager {
	return &Manager{
		rollback: func() error {
			return nil
		},
		commit: func() error {
			return nil
		},
	}
}
