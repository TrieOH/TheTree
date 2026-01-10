package apierr

var (
	ErrInvalidInput = Error{
		Code:    InvalidInput,
		Message: "invalid input",
		ID:      SystemInternalError,
	}

	ErrBadRequest = Error{
		Code:    BadRequest,
		Message: "bad request",
		ID:      SystemInternalError,
	}

	ErrNotFound = Error{
		Code:    NotFound,
		Message: "resource not found",
		ID:      SystemInternalError,
	}

	ErrConflict = Error{
		Code:    Conflict,
		Message: "conflict",
		ID:      SystemInternalError,
	}

	ErrUnauthorized = Error{
		Code:    Unauthorized,
		Message: "unauthorized",
		ID:      SystemInternalError,
	}

	ErrForbidden = Error{
		Code:    Forbidden,
		Message: "forbidden",
		ID:      SystemInternalError,
	}

	ErrInternal = Error{
		Code:    Internal,
		Message: "internal error",
		ID:      SystemInternalError,
	}
)
