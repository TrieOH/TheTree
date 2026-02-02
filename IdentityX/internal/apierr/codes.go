package apierr

import "errors"

const (
	InvalidInput Code = "INVALID_INPUT"
	BadRequest   Code = "BAD_REQUEST"
	NotFound     Code = "NOT_FOUND"
	Conflict     Code = "CONFLICT"
	Unauthorized Code = "UNAUTHORIZED"
	Forbidden    Code = "FORBIDDEN"
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
