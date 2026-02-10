package liblog

import (
	"errors"
)

// Level is logging level.
type Level string

const (
	// TraceLevel is the value used for the trace logging level.
	TraceLevel Level = "trace"
	// DebugLevel is the value used for the debug logging level.
	DebugLevel Level = "debug"
	// InfoLevel is the value used for the info logging level.
	InfoLevel Level = "info"
	// WarnLevel is the value used for the warning logging level.
	WarnLevel Level = "warning"
	// ErrorLevel is the value used for the error logging level.
	ErrorLevel Level = "error"
	// FatalLevel is the value used for the fatal logging level.
	FatalLevel Level = "fatal"
	// PanicLevel is the value used for the panic logging level.
	PanicLevel Level = "panic"
)

// Validate validates Level logging value.
func (l Level) Validate() error {
	if l == "" {
		return errEmptyLogLevel
	}

	switch l {
	case
		PanicLevel,
		FatalLevel,
		ErrorLevel,
		WarnLevel,
		InfoLevel,
		DebugLevel,
		TraceLevel:
		return nil
	default:
		return errInvalidLogLevel
	}
}

// Output is writer stream.
type Output string

const (
	// StdOut is standard output.
	StdOut Output = "StdOut"
	// StdErr is standard error output.
	StdErr Output = "StdErr"
	// ConsoleOutput is standard output for local debugging to enable colors.
	ConsoleOutput Output = "ConsoleStdOutput"
)

// Validate validates the output stream.
func (o Output) Validate() error {
	if o == "" {
		return errEmptyLoggerOutput
	}

	return nil
}

// Outputs is slice of Output.
type Outputs []Output

// Validate validates the outputs stream.
func (o Outputs) Validate() error {
	if len(o) == 0 {
		return errEmptyLoggerOutput
	}

	mapListOfOutputs := make(map[Output]struct{}, len(o))

	for _, output := range o {
		// checking for uniqueness
		if _, exist := mapListOfOutputs[output]; exist {
			return errDuplicateLogOutput
		}

		mapListOfOutputs[output] = struct{}{}

		// validating output
		if err := output.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Format is.
type Format string

const (
	// FormatJSON is json log format.
	FormatJSON = "json"
	// FormatText is text log format.
	FormatText = "text"
	// FormatHumanReadable is human-readable log format.
	FormatHumanReadable = "human-readable"
)

// Validate validates log format.
func (f Format) Validate() error {
	if f == "" {
		return errEmptyLogFormat
	}

	switch f {
	case
		FormatJSON,
		FormatText,
		FormatHumanReadable:
		return nil
	default:
		return errInvalidLogFormat
	}
}

// Config is logger configs.
type Config struct {
	Level           Level        `env:"LOGGER_LEVEL" json:"level" yaml:"level"`
	Outputs         Outputs      `env:"LOGGER_OUTPUTS" json:"outputs" yaml:"outputs"`
	Formatter       Format       `env:"LOGGER_FORMATTER" json:"formatter" yaml:"formatter"`
	TimeStampFormat string       `env:"LOGGER_TIMESTAMP_FORMAT" json:"time_stamp_format" yaml:"time_stamp_format"`
	Caller          bool         `env:"LOGGER_CALLER" json:"caller" yaml:"caller"`
	Sentry          SentryConfig `json:"sentry" yaml:"sentry"`
}

// SentryConfig is the configuration of the Sentry.
type SentryConfig struct {
	DSN         string `env:"LOGGER_SENTRY_DSN" json:"dsn" yaml:"dns"`
	Environment string `env:"LOGGER_SENTRY_ENV" json:"environment" yaml:"environment"`
	Release     string `env:"LOGGER_SENTRY_RELEASE" json:"release" yaml:"release"`
}

// Validate validates logger configs.
//
//nolint:revive // here we will have new validations so if-return check will disturb us
func (c Config) Validate() error {
	if err := c.Level.Validate(); err != nil {
		return err
	}

	if err := c.Outputs.Validate(); err != nil {
		return err
	}

	if err := c.Formatter.Validate(); err != nil {
		return err
	}

	return nil
}

var (
	errEmptyLoggerOutput  = errors.New("empty log output stream specified")
	errInvalidLogLevel    = errors.New("invalid log level specified")
	errEmptyLogLevel      = errors.New("empty log level specified")
	errInvalidLogFormat   = errors.New("invalid log format specified")
	errEmptyLogFormat     = errors.New("empty log format specified")
	errDuplicateLogOutput = errors.New("duplicate log output specified")
)
