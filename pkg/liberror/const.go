package liberror

import "net/http"

var (
	// ErrInternal is returned when there is an internal server error.
	ErrInternal = &Error{
		Err:      "internal server error",
		Code:     "INTERNAL_ERROR",
		HTTPCode: http.StatusInternalServerError,
	}

	// ErrBadRequest is returned when there is no valid request.
	ErrBadRequest = &Error{
		Err:      "bad request error",
		Code:     "BAD_REQUEST",
		HTTPCode: http.StatusBadRequest,
	}

	// ErrNotFound is returned when the resource is not found.
	ErrNotFound = &Error{
		Err:      "not found error",
		Code:     "NOT_FOUND",
		HTTPCode: http.StatusNotFound,
	}
)
