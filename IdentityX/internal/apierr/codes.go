package apierr

import "errors"

const (
	InvalidInput Code = "INVALID_INPUT"
	NotFound     Code = "NOT_FOUND"
	Conflict     Code = "CONFLICT"
	Unauthorized Code = "UNAUTHORIZED"
	Forbidden    Code = "FORBIDDEN"
	Internal     Code = "INTERNAL"
)

func As(err error) (*Error, bool) {
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

func IsNotFound(err error) bool {
	apiErr, ok := As(err)
	return ok && apiErr.Code == NotFound
}

func IsInvalidInput(err error) bool {
	apiErr, ok := As(err)
	return ok && apiErr.Code == InvalidInput
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

func IsSystemError(err error) bool {
	apiErr, ok := As(err)
	if !ok {
		return true // unknown error = system
	}

	switch apiErr.Code {
	case InvalidInput, NotFound, Conflict, Unauthorized, Forbidden:
		return false
	default:
		return true
	}
}
