// Package postgres is a driver to work with postgresql database.
package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DB is a client for interacting with the PostgresSQL database.
type DB interface {
	NamedQuery(query string, arg any) (*sqlx.Rows, error)
	QueryRowx(query string, args ...any) *sqlx.Row

	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)

	Ping() error
}
