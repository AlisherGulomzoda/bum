//nolint:ireturn // we are writing wrapper over two different structures, and we cannot avoid using interface.
package transaction

import (
	"context"
	"database/sql"
	"fmt"

	"bum-service/pkg/postgres"
)

type transactionKey struct{}

// Injector надстройка, которая позволяет открывать транзакцию,
// передавая контроль над результатом стороннему потребителю,
// передавая транзакцию через контекст.
type Injector struct {
	db postgres.DB
}

// NewInjector создает новый экземпляр Injector.
func NewInjector(db postgres.DB) *Injector {
	return &Injector{
		db: db,
	}
}

// Inject начинает транзакцию и запечатывает ее в контекст.
func (c Injector) Inject(ctx context.Context) (context.Context, *Manager, error) {
	if _, ok := c.ExtractDB(ctx); ok {
		return ctx, noopSolver(), nil
	}

	pgxTx, err := c.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	client := postgres.SqlxTxClient{
		Tx: pgxTx,
	}

	ctx = context.WithValue(ctx, transactionKey{}, client)

	return ctx, newSolver(client), nil
}

// ExtractDB возвращает транзакцию из контекста или создает новую если
// в контексте транзакция не найдена.
func (c Injector) ExtractDB(ctx context.Context) (db postgres.DB, isTx bool) {
	tx, ok := ctx.Value(transactionKey{}).(postgres.SqlxTxClient)
	if !ok {
		return c.db, false
	}

	return tx, ok
}
