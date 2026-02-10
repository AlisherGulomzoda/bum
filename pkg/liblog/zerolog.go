//nolint:nolintlint,ireturn // in order to implement the Logger interface we have to return this interface
package liblog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	zero "github.com/rs/zerolog"
)

// ZeroLog is a wrapper over zero-log.
type ZeroLog struct {
	*zero.Logger
}

const (
	colorBold = 1
)

// ZeroLogOption represents optional settings for logger.
type ZeroLogOption func(op *ZeroLog) error

// NewZeroLog returns a new ZeroLog.
func NewZeroLog(conf Config, options ...ZeroLogOption) (*ZeroLog, error) {
	// setting up our log level
	level, err := levelToZeroLogLevel(conf.Level)
	if err != nil {
		return nil, err
	}

	zero.SetGlobalLevel(level)

	// getting log writers
	writer, err := getWriter(conf)
	if err != nil {
		return nil, err
	}

	// creating a new logger
	logger := zero.New(
		writer,
	).With().Timestamp().Logger()

	// setting up caller
	logger = setCaller(conf, logger)

	// setting up time format
	setZeroLogTimeStamp(conf)

	// adding Process ID, Parent Process ID and Hostname to our logger
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	logger = logger.With().Fields(map[string]any{
		"pid":      os.Getpid(),
		"ppid":     os.Getppid(),
		"hostname": hostname,
	}).Logger()

	zeroLog := &ZeroLog{
		Logger: &logger,
	}

	// setting options
	for _, op := range options {
		if err = op(zeroLog); err != nil {
			return nil, err
		}
	}

	return zeroLog, nil
}

// getWriter returns a writer for the given configuration.
func getWriter(conf Config) (io.Writer, error) {
	var writer io.Writer

	writers := make([]io.Writer, 0, len(conf.Outputs))

	for _, cfgOutput := range conf.Outputs {
		var w io.Writer

		w, err := cfgOutputToZeroLogOutput(cfgOutput)
		if err != nil {
			return nil, err
		}

		// setting up log format
		if w, err = setZeroLogFormat(
			w,
			conf.Formatter,
			isZeroLogColorDisabled(cfgOutput),
		); err != nil {
			return nil, err
		}

		writers = append(writers, w)
	}

	if len(conf.Outputs) > 1 {
		writer = zero.MultiLevelWriter(writers...)
	}

	return writer, nil
}

// setCaller sets caller to our zero logger.
func setCaller(conf Config, logger zero.Logger) zero.Logger {
	if conf.Caller {
		// setting up "Skip Frame" because we use wrapper over zero log
		return logger.With().CallerWithSkipFrameCount(SkipFrames).Logger()
	}

	return logger
}

// isZeroLogColorDisabled returns true if it is local debugging mode.
func isZeroLogColorDisabled(output Output) bool {
	return output != ConsoleOutput
}

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s any, c int, isColorDisabled bool) string {
	if isColorDisabled {
		return fmt.Sprintf("%s", s)
	}

	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// levelToZeroLogLevel converts Level to zero log level.
func levelToZeroLogLevel(level Level) (zero.Level, error) {
	switch level {
	case TraceLevel:
		return zero.TraceLevel, nil
	case DebugLevel:
		return zero.DebugLevel, nil
	case InfoLevel:
		return zero.InfoLevel, nil
	case WarnLevel:
		return zero.WarnLevel, nil
	case ErrorLevel:
		return zero.ErrorLevel, nil
	case FatalLevel:
		return zero.FatalLevel, nil
	case PanicLevel:
		return zero.PanicLevel, nil
	default:
		return -1, errInvalidZeroLogLevel
	}
}

const (
	defaultLogFileMode os.FileMode = 0o644
)

// cfgOutputToZeroLogOutput converts conf output to zero output.
func cfgOutputToZeroLogOutput(confOutput Output) (io.Writer, error) {
	var output io.Writer

	switch confOutput {
	case StdOut:
		return os.Stdout, nil
	case StdErr:
		return os.Stderr, nil
	case ConsoleOutput:
		return os.Stdout, nil
	default:
		f, err := os.OpenFile(string(confOutput), os.O_APPEND|os.O_CREATE|os.O_WRONLY, defaultLogFileMode)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file [%s] : %w", string(confOutput), err)
		}

		output = f
	}

	return output, nil
}

// setZeroLogFormat sets zero log format.
func setZeroLogFormat(
	writer io.Writer,
	format Format,
	isColorDisabled bool,
) (io.Writer, error) {
	switch format {
	case FormatJSON:
		// by default zero log format is json
		return writer, nil
	case FormatHumanReadable:
		return getHumanReadableWriter(writer, isColorDisabled)
	default:
		return nil, errInvalidZeroLogFormat
	}
}

func getHumanReadableWriter(writer io.Writer, isColorDisabled bool) (io.Writer, error) {
	return zero.ConsoleWriter{
		Out:     writer,
		NoColor: isColorDisabled,
		FormatTimestamp: func(i any) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatCaller: func(i any) string {
			var c string
			if cc, ok := i.(string); ok {
				c = cc
			}
			if c != "" {
				if cwd, err := os.Getwd(); err == nil {
					if rel, err := filepath.Rel(cwd, c); err == nil {
						c = rel
					}
				}
				c = colorize(c, colorBold, isColorDisabled)
			}

			return c
		},
	}, nil
}

// setZeroLogTimeStamp sets zero log timestamp format.
func setZeroLogTimeStamp(conf Config) {
	if conf.TimeStampFormat == "" {
		return
	}

	zero.TimeFieldFormat = conf.TimeStampFormat //nolint:reassign // this is time format setting for zero log library
}

// Fatalf starts a new fatal formatted message.
func (z ZeroLog) Fatalf(format string, args ...any) {
	z.Logger.Fatal().Msgf(format, args...)
}

// Fatal starts a new error message.
func (z ZeroLog) Fatal(args ...any) {
	z.Logger.Fatal().Msg(fmt.Sprint(args...))
}

// Errorf starts a new error formatted message.
func (z ZeroLog) Errorf(format string, args ...any) {
	z.Logger.Error().Msgf(format, args...)
}

// Error starts a new error message.
func (z ZeroLog) Error(args ...any) {
	z.Logger.Error().Msg(fmt.Sprint(args...))
}

// Warningf starts a new Warn formatted message.
func (z ZeroLog) Warningf(format string, args ...any) {
	z.Logger.Warn().Msgf(format, args...)
}

// Warning starts a new Warn message.
func (z ZeroLog) Warning(args ...any) {
	z.Logger.Warn().Msg(fmt.Sprint(args...))
}

// Infof starts a new info formatted message.
func (z ZeroLog) Infof(format string, args ...any) {
	z.Logger.Info().Msgf(format, args...)
}

// Info starts a new info message.
func (z ZeroLog) Info(args ...any) {
	z.Logger.Info().Msg(fmt.Sprint(args...))
}

// Debugf starts a new debug formatted message.
func (z ZeroLog) Debugf(format string, args ...any) {
	z.Logger.Debug().Msgf(format, args...)
}

// Debug starts a new debug message.
func (z ZeroLog) Debug(args ...any) {
	z.Logger.Debug().Msg(fmt.Sprint(args...))
}

// Tracef starts a new trace formatted message.
func (z ZeroLog) Tracef(format string, args ...any) {
	z.Logger.Trace().Msgf(format, args...)
}

// Trace starts a new trace message.
func (z ZeroLog) Trace(args ...any) {
	z.Logger.Trace().Msg(fmt.Sprint(args...))
}

// WithLevel sets level to use for further Msg method call.
//
//nolint:errcheck,ireturn // We don't want to check because we just need to print
func (z ZeroLog) WithLevel(level Level) LogFinalizer {
	zeroLevel, _ := levelToZeroLogLevel(level)

	return z.Logger.WithLevel(zeroLevel)
}

// WithFields adds some fields to the logger.
//
//nolint:nolintlint,ireturn // in order to implement the Logger interface we have to return this interface
func (z ZeroLog) WithFields(v map[string]any) Logger {
	logWithFields := z.With().Fields(v).Logger()

	return ZeroLog{
		Logger: &logWithFields,
	}
}

// CallerWithSkipFrameCount sets the skip frame count and returns the new logger.
func (z ZeroLog) CallerWithSkipFrameCount(skipFrameCount int) Logger {
	logWithSkipFrame := z.Logger.With().CallerWithSkipFrameCount(skipFrameCount).Logger()

	return ZeroLog{
		Logger: &logWithSkipFrame,
	}
}

var (
	errInvalidZeroLogLevel  = errors.New("invalid zero log level value")
	errInvalidZeroLogFormat = errors.New("invalid zero log format")
)
