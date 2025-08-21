package customerrors

import (
	"fmt"
	"net/http"
)

type Error struct {
	message string
	code    int
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Code() int {
	return e.code
}

var (
	ERROR_UNAUTHORISED              = &Error{message: "unauthorised", code: http.StatusUnauthorized}
	ERROR_DATABASE                  = &Error{message: "database operation failed", code: http.StatusFailedDependency}
	ERROR_PERMISSION_DENIED         = &Error{message: "permission denied", code: http.StatusForbidden}
	ERROR_DATABASE_UNIQUE_KEY       = &Error{message: "unique key failed", code: http.StatusBadRequest}
	ERROR_DATABASE_RECORD_NOT_FOUND = &Error{message: "record not found", code: http.StatusNotFound}
	ERROR_INTERNAL_SERVER_ERROR     = &Error{message: "internal server error", code: http.StatusInternalServerError}

	// Database connection errors
	ERROR_DATABASE_CONNECTION_FAILED = &Error{message: "failed to connect to database", code: http.StatusInternalServerError}
	// PING FAILED
	ERROR_DATABASE_PING_FAILED = &Error{message: "failed to ping database", code: http.StatusInternalServerError}
)

func BAD_REQUEST_ERROR(message string, args ...any) error {
	return &Error{code: http.StatusBadRequest, message: fmt.Sprintf(message, args...)}
}

func New(code int, message string, args ...any) error {
	return &Error{code: code, message: fmt.Sprintf(message, args...)}
}
