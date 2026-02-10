package liblog

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	// SkipFrames is frame for logging.
	SkipFrames = 3
)

// Fields is a map of field names to add to our logger.
type Fields map[string]any

// Logger is common logger interface.
type Logger interface {
	Errorf(format string, args ...any)
	Error(args ...any)
	Warningf(format string, args ...any)
	Warning(args ...any)
	Infof(format string, args ...any)
	Info(args ...any)
	Debugf(format string, args ...any)
	Debug(args ...any)
	Tracef(format string, args ...any)
	Trace(args ...any)
	WithLevel(level Level) LogFinalizer
	WithFields(v map[string]any) Logger
	CallerWithSkipFrameCount(skipFrameCount int) Logger
}

// LogFinalizer is used for logging after specified logger level.
type LogFinalizer interface {
	Msg(args string)
}

type contextKey struct{}

// With - returns a copy of parent in which the value associated with key is a logger.
func With(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

// Must - extracts a logger instance from context panicking on failure.
//
//nolint:nolintlint,forcetypeassert,ireturn // it's must method that was intended to panic
func Must(ctx context.Context) Logger {
	return ctx.Value(contextKey{}).(Logger)
}

// From gets logger from the provided context.
//
//nolint:nolintlint,ireturn // or logger is an interface so we return this interface
func From(ctx context.Context) (Logger, bool) {
	logger, ok := ctx.Value(contextKey{}).(Logger)
	return logger, ok
}

// WithRequest sets context with data from request.
func WithRequest(ctx context.Context, r *http.Request) context.Context {
	entry := WithRequestFields(Must(ctx), r)
	return With(ctx, entry)
}

const (
	// ClientIPHeader is client IP header.
	ClientIPHeader = "X-Forwarded-For"
	// UserAgentHeader is user agent header.
	UserAgentHeader = "User-Agent"
	// SessionIDHeader is Authentication token header.
	SessionIDHeader = "X-Authentication-Token"
	// ReqIDHeader is request ID header.
	ReqIDHeader = "X-Request-ID"
)

// WithRequestFields sets logger fields with data from request.
//
//nolint:ireturn // in order to implement our Logger interface we need to return interface
func WithRequestFields(logger Logger, r *http.Request) Logger {
	return logger.WithFields(
		Fields{
			"client_ip":      r.Header.Get(ClientIPHeader),
			"user_agent":     r.Header.Get(UserAgentHeader),
			"session_id":     r.Header.Get(SessionIDHeader),
			"request_id":     r.Header.Get(ReqIDHeader),
			"request_path":   r.URL.Path,
			"request_query":  r.URL.RawQuery,
			"request_uri":    r.URL.RequestURI(),
			"request_method": r.Method,
		},
	)
}

// NewMiddleware creates a middleware which contains common
// actions for all services like logger settings, tracing, etc.
func NewMiddleware(l Logger) func(handler http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				const header = "X-Request-Id"
				if r.Header.Get(header) == "" {
					r.Header.Set(header, uuid.New().String())
				}

				ctx := WithRequest(
					With(r.Context(), l),
					r,
				)

				h.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
