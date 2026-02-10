package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	_defaultMaxOpenConnections     = 25
	_defaultMaxIdleConnections     = 10
	_defaultMaxLifetimeConnections = time.Minute

	_defaultMaxAttempts     = 3
	_defaultAttemptInterval = 5 * time.Second

	pgxDriver = "pgx"
)

// Option is optional argument.
type Option func(*SqlxDBClient) error

// SqlxDBClient is a pgx pool client.
type SqlxDBClient struct {
	maxOpenConnections     int
	maxIdleConnections     int
	maxLifetimeConnections time.Duration

	*sqlx.DB
}

// SqlxTxClient is a pgx Tx client.
type SqlxTxClient struct {
	*sqlx.Tx
}

// BeginTxx returns nil.
//
//nolint:nilnil // it's ok
func (SqlxTxClient) BeginTxx(_ context.Context, _ *sql.TxOptions) (*sqlx.Tx, error) {
	return nil, nil
}

// Ping pings the database connection.
func (SqlxTxClient) Ping() error {
	return nil
}

// NewClient creates a new pgx pool client.
func NewClient(_ context.Context, url string, opts ...Option) (*SqlxDBClient, error) {
	var (
		pgxClient = &SqlxDBClient{
			maxOpenConnections:     _defaultMaxOpenConnections,
			maxIdleConnections:     _defaultMaxIdleConnections,
			maxLifetimeConnections: _defaultMaxLifetimeConnections,
		}
		err error
	)

	// Custom options
	for _, opt := range opts {
		if err = opt(pgxClient); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// we will try a few times to connect to database
	for attempt := 1; ; attempt++ {
		if attempt > _defaultMaxAttempts {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		pgxClient.DB, err = sqlx.Connect(pgxDriver, url)
		if err != nil {
			time.Sleep(_defaultAttemptInterval)
			continue
		}

		if err = pgxClient.Ping(); err != nil {
			time.Sleep(_defaultAttemptInterval)
			continue
		}

		break
	}

	pgxClient.SetMaxOpenConns(pgxClient.maxOpenConnections)
	pgxClient.SetMaxIdleConns(pgxClient.maxIdleConnections)
	pgxClient.SetConnMaxLifetime(pgxClient.maxLifetimeConnections)

	return pgxClient, nil
}

// Close closes the connection.
func (p SqlxDBClient) Close(_ context.Context) error {
	if err := p.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// WithMaxOpenConnections sets the maximum number of open connections to the database.
func WithMaxOpenConnections(size int) Option {
	return func(client *SqlxDBClient) error {
		client.maxOpenConnections = size

		return nil
	}
}

// WithMaxIdleConnections sets the maximum number of connections in the idle connection pool.
func WithMaxIdleConnections(size int) Option {
	return func(client *SqlxDBClient) error {
		client.maxIdleConnections = size

		return nil
	}
}

// WithMaxLifetimeConnections sets the maximum amount of time a connection may be reused.
func WithMaxLifetimeConnections(d time.Duration) Option {
	return func(client *SqlxDBClient) error {
		client.maxLifetimeConnections = d

		return nil
	}
}
