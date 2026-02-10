package liberror

// WrapOption is optional parameter for Wrap function.
type WrapOption func(err *Error)

// WithHTTPCode returns WrapOption which set HTTPCode field of base to err.
func WithHTTPCode(code int) WrapOption {
	return func(err *Error) {
		err.HTTPCode = code
	}
}

// WithErr returns WrapOption which set Err field of base to err.
func WithErr(errStr string) WrapOption {
	return func(err *Error) {
		err.Err = errStr
	}
}

// Wrap wraps Error structure.
// ErrInternal is used as a base error.
// Err becomes child of base error.
// WrapOptions evaluates on base error before the child is installed.
func Wrap(err *Error, opts ...WrapOption) *Error {
	base := *err

	for _, opt := range opts {
		opt(&base)
	}

	base.child = err

	return &base
}
