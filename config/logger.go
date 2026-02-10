package config

// Logger is logger configs.
type Logger struct {
	Level           string       `yaml:"level" validate:"required"`
	Outputs         []string     `yaml:"outputs" validate:"required"`
	Formatter       string       `yaml:"formatter" validate:"required"`
	TimeStampFormat string       `yaml:"time_stamp_format" validate:"required"`
	Caller          bool         `yaml:"caller" validate:"required"`
	Sentry          SentryConfig `yaml:"sentry"`
}

// SentryConfig is the configuration of the Sentry.
type SentryConfig struct {
	DSN         string `yaml:"dns"`
	Environment string `yaml:"environment"`
	Release     string `yaml:"release"`
}
