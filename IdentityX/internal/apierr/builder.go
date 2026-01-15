package apierr

func (e Error) WithMsg(msg string) *Error {
	e.Message = msg
	return &e
}

// WithCause appends an error to the Causes slice.
// For backwards compatibility, it also sets the legacy Cause field to the first error.
func (e Error) WithCause(err error) *Error {
	if err == nil {
		return &e
	}

	// Initialize Causes if needed
	if e.Causes == nil {
		e.Causes = make([]error, 0, 1)
	}

	// Append to causes
	e.Causes = append(e.Causes, err)

	// Backwards compatibility: set Cause to first error
	if e.Cause == nil {
		e.Cause = err
	}

	return &e
}

// WithDebugCause appends an error to the DebugCauses slice.
// These causes are intended for debugging and can be safely disabled in production.
func (e Error) WithDebugCause(err error) *Error {
	if err == nil {
		return &e
	}

	// Initialize DebugCauses if needed
	if e.DebugCauses == nil {
		e.DebugCauses = make([]error, 0, 1)
	}

	e.DebugCauses = append(e.DebugCauses, err)
	return &e
}

func (e Error) WithCode(code Code) *Error {
	e.Code = code
	return &e
}

func (e Error) WithID(id ID) *Error {
	e.ID = id
	return &e
}

func (e Error) WithField(key string, value any) *Error {
	if e.Fields == nil {
		e.Fields = map[string]any{}
	}
	e.Fields[key] = value
	return &e
}
