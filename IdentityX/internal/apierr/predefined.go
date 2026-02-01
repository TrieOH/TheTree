package apierr

var (
	ErrInvalidInput = Error{
		Code:    InvalidInput,
		Message: "invalid input",
		ID:      PlaceholderID,
	}

	ErrBadRequest = Error{
		Code:    BadRequest,
		Message: "bad request",
		ID:      PlaceholderID,
	}

	ErrNotFound = Error{
		Code:    NotFound,
		Message: "resource not found",
		ID:      PlaceholderID,
	}

	ErrConflict = Error{
		Code:    Conflict,
		Message: "conflict",
		ID:      PlaceholderID,
	}

	ErrUnauthorized = Error{
		Code:    Unauthorized,
		Message: "unauthorized",
		ID:      PlaceholderID,
	}

	ErrForbidden = Error{
		Code:    Forbidden,
		Message: "forbidden",
		ID:      PlaceholderID,
	}

	ErrInternal = Error{
		Code:    Internal,
		Message: "internal error",
		ID:      PlaceholderID,
	}
)
