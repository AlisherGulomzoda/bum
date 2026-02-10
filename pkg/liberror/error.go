package liberror

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Error represents json error with http code and error.
type Error struct {
	Err      string `json:"error,omitempty"`
	Code     string `json:"code,omitempty"`
	HTTPCode int    `json:"-"`
	child    error
}

// MarshalJSON is an implementation of the MarshalJSON interface in encoding/json.
func (e *Error) MarshalJSON() ([]byte, error) {
	type t struct {
		Err      string `json:"error,omitempty"`
		Code     string `json:"code,omitempty"`
		HTTPCode int    `json:"-"`
		child    error
	}

	// Casting is needed because json.Marshal(e) creates infinite recursion
	// to the MarshalJSON method of CustomError.
	b, err := json.Marshal((*t)(e))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal liberror : %w", err)
	}

	return b, nil
}

// StatusCode is an implementation of the StatusCoder interface.
func (e *Error) StatusCode() int {
	return e.HTTPCode
}

// Error is an implementation of the Error interface.
//
//nolint:revive // we don't need to handle errors here
func (e *Error) Error() string {
	errorBuilder := strings.Builder{}
	errorBuilder.WriteString(e.Err)

	if e.child != nil {
		errorBuilder.WriteString(fmt.Sprintf(";\nChild = [%s]", e.child))
	}

	return errorBuilder.String()
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (e *Error) Unwrap() error {
	return e.child
}

// Is reports whether any error in err's chain matches target.
//
//nolint:errorlint // we are implementing interface
func (e *Error) Is(target error) bool {
	if e2, ok := target.(*Error); ok {
		return e.Err == e2.Err &&
			e.Code == e2.Code &&
			e.HTTPCode == e2.HTTPCode
	}

	return false
}
