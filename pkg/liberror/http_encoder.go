package liberror

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// StatusCoder is checked by ErrorEncoder.
// If an error value implements StatusCoder, the StatusCode will be used when
// encoding the error. By default, StatusInternalServerError (500) is used.
type StatusCoder interface {
	StatusCode() int
}

// ErrorEncoder is responsible for encoding an error to the ResponseWriter.
// For errors without json.Marshaler returns ErrInternal error.
//
//nolint:errorlint // we don't check for specific error, and instead we want to find out does it have StatusCode method
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var (
		code = http.StatusInternalServerError

		response json.Marshaler = ErrInternal

		customErr *Error
	)

	if errors.As(err, &customErr) {
		code = customErr.HTTPCode
		response = customErr
	}

	w.WriteHeader(code)

	if errEncode := json.NewEncoder(w).Encode(response); errEncode != nil {
		return fmt.Errorf("failed to encode: %w", errEncode)
	}

	return nil
}
