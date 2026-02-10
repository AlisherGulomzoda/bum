//nolint:ireturn // we are writing wrapper over two different structures, and we cannot avoid using interface.
package transaction

import (
	"context"
	"fmt"

	"bum-service/pkg/postgres"
)

// Session управляет запуском сессии на application уровне.
type Session interface {
	Begin(context.Context) (ctx context.Context, solver SessionSolver, err error)
	End(tx SessionSolver, err error) error
}

// SessionSolver управляет результатом транзакции в рамках сессии.
type SessionSolver interface {
	Rollback() error
	Commit() error
}

// SessionDB управляет доступом к сессионному соединению с БД.
type SessionDB interface {
	DB(context.Context) postgres.DB
}

// SessionAdapter реализует интерфейс Session и инкапсулирует взаимодействие
// с зависимостями, зависящими от сессии.
type SessionAdapter struct {
	injector *Injector
}

// NewSessionAdapter создает новый инстанс SessionAdapter.
func NewSessionAdapter(db postgres.DB) *SessionAdapter {
	return &SessionAdapter{
		injector: NewInjector(db),
	}
}

// Begin запускает сессию для зависимых от сессии объектов.
func (s SessionAdapter) Begin(
	ctx context.Context,
) (context.Context, SessionSolver, error) {
	return s.injector.Inject(ctx)
}

// End завершает сессию для фиксации изменений или откат изменений в зависимости от передаваемой ошибки.
func (SessionAdapter) End(tx SessionSolver, err error) error {
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return fmt.Errorf("failed to rollback transaction, error %w : %w",
				err, errRollback)
		}

		return nil
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return fmt.Errorf("failed to commit transaction: %w", errCommit)
	}

	return nil
}

// DB возвращает соединение с БД из сессии.
func (s SessionAdapter) DB(ctx context.Context) postgres.DB {
	db, _ := s.injector.ExtractDB(ctx)
	return db
}
