package liblog

// DummyLogger is a dummy logger.
type DummyLogger struct{}

// CallerWithSkipFrameCount is a dummy CallerWithSkipFrameCount method.
//
//nolint:ireturn // it's an implementation of logger interface.
func (z DummyLogger) CallerWithSkipFrameCount(_ int) Logger { return z }

// NewDummyLogger returns a new DummyLogger.
func NewDummyLogger() *DummyLogger { return &DummyLogger{} }

// DummyFinalizer is dummy finalizer.
type DummyFinalizer struct{}

// Msg starts a new message.
func (DummyFinalizer) Msg(_ string) {}

// Fatalf starts a new fatal formatted message.
func (DummyLogger) Fatalf(_ string, _ ...any) {}

// Fatal starts a new error message.
func (DummyLogger) Fatal(_ ...any) {}

// Errorf starts a new error formatted message.
func (DummyLogger) Errorf(_ string, _ ...any) {}

// Error starts a new error message.
func (DummyLogger) Error(_ ...any) {}

// Warningf starts a new Warn formatted message.
func (DummyLogger) Warningf(_ string, _ ...any) {}

// Warning starts a new Warn message.
func (DummyLogger) Warning(_ ...any) {}

// Infof starts a new info formatted message.
func (DummyLogger) Infof(_ string, _ ...any) {}

// Info starts a new info message.
func (DummyLogger) Info(_ ...any) {}

// Debugf starts a new debug formatted message.
func (DummyLogger) Debugf(_ string, _ ...any) {}

// Debug starts a new debug message.
func (DummyLogger) Debug(_ ...any) {}

// Tracef starts a new trace formatted message.
func (DummyLogger) Tracef(_ string, _ ...any) {}

// Trace starts a new trace message.
func (DummyLogger) Trace(_ ...any) {}

// WithLevel sets level to use for further Msg method call.
//
//nolint:errcheck,ireturn // We don't want to check because we just need to print
func (DummyLogger) WithLevel(_ Level) LogFinalizer { return DummyFinalizer{} }

// WithFields adds some fields to the logger.
//
//nolint:nolintlint,ireturn // in order to implement the Logger interface we have to return this interface
func (z DummyLogger) WithFields(_ map[string]any) Logger { return z }
