package config

import (
	"fmt"
	"time"
)

// Database is a database configuration structure.
type Database struct {
	// SQL-Dialect
	Dialect                         string        `env:"DB_DIALECT" yaml:"dialect" validate:"required"`
	User                            string        `env:"DB_USER" yaml:"user" validate:"required"`
	Password                        string        `env:"DB_PASSWORD" yaml:"password" validate:"required"`
	Name                            string        `env:"DB_NAME" yaml:"name" validate:"required"`
	Host                            string        `env:"DB_HOST" yaml:"host" validate:"required"`
	Port                            int           `env:"DB_PORT" yaml:"port" validate:"required"`
	SSLMode                         string        `env:"DB_SSL_MODE" yaml:"ssl_mode" validate:"required"`
	Schema                          string        `env:"DB_SCHEMA" yaml:"schema" validate:"required"`
	IdleInTransactionSessionTimeout time.Duration `env:"DB_IDLE_IN_TRANSACTION_SESSION_TIMEOUT" yaml:"idle_in_transaction_session_timeout" validate:"required"` //nolint:lll // This is a tag
}

// GetDatabaseDSN returns DSN from database configuration.
func (db Database) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s search_path=%s",
		db.Host,
		db.Port,
		db.User,
		db.Name,
		db.Password,
		db.SSLMode,
		db.Schema,
	)
}
