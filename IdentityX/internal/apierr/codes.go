package apierr

import "errors"

const (
	InvalidInput Code = "INVALID_INPUT"
	BadRequest   Code = "BAD_REQUEST"
	NotFound     Code = "NOT_FOUND"
	Conflict     Code = "CONFLICT"
	Unauthorized Code = "UNAUTHORIZED"
	Forbidden    Code = "FORBIDDEN"
	Internal     Code = "INTERNAL"
)

// As converts an error to an API error.
// It returns the API error and true if the error is an API error, otherwise it returns nil and false.
func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

func IsConflict(err error) bool {
	apiErr, ok := As(err)
	return ok && apiErr.Code == Conflict
}

func IsInvalidInput(err error) bool {
	apiErr, ok := As(err)
	return ok && apiErr.Code == InvalidInput
}

func IsBadRequest(err error) bool {
	apiErr, ok := As(err)
	return ok && apiErr.Code == BadRequest
}

func IsUnauthorized(err error) bool {
	apiErr, ok := As(err)
	return ok && apiErr.Code == Unauthorized
}

func IsForbidden(err error) bool {
	apiErr, ok := As(err)
	return ok && apiErr.Code == Forbidden
}

func IsInternal(err error) bool {
	apiErr, ok := As(err)
	return ok && apiErr.Code == Internal
}

// IsSystemError returns true if the error is a system error.
// An error is considered a system error if it is not an API error or if it is an internal API error.
func IsSystemError(err error) bool {
	apiErr, ok := As(err)
	if !ok {
		return true // unknown error = system
	}

	switch apiErr.Code {
	case InvalidInput, BadRequest, NotFound, Conflict, Unauthorized, Forbidden:
		return false
	default:
		return true
	}
}
