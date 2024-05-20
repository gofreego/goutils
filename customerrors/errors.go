package customerrors

import (
	"fmt"
	"net/http"
)

type Error struct {
	Message string
	Code    int
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ERROR_UNAUTHORISED              = &Error{Message: "unauthorised", Code: http.StatusUnauthorized}
	ERROR_DATABASE                  = &Error{Message: "database operation failed", Code: http.StatusFailedDependency}
	ERROR_PERMISSION_DENIED         = &Error{Message: "permission denied", Code: http.StatusForbidden}
	ERROR_DATABASE_UNIQUE_KEY       = &Error{Message: "unique key failed", Code: http.StatusBadRequest}
	ERROR_DATABASE_RECORD_NOT_FOUND = &Error{Message: "record not found", Code: http.StatusNotFound}
	ERROR_INTERNAL_SERVER_ERROR     = &Error{Message: "internal server error", Code: http.StatusInternalServerError}
)

func BAD_REQUEST_ERROR(message string, args ...any) *Error {
	return &Error{Code: http.StatusBadRequest, Message: fmt.Sprintf(message, args...)}
}
