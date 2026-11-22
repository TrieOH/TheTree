package errs

// TemplateRenderError is an email template override that does not render.
type TemplateRenderError struct{ Err error }

func (e *TemplateRenderError) Error() string {
	return "email template does not render: " + e.Err.Error()
}

func (e *TemplateRenderError) Unwrap() error { return e.Err }

// ProfileSchemaViolationError is a profile its scope's schema rejects; the
// message is the validator's.
type ProfileSchemaViolationError struct{ Err error }

func (e *ProfileSchemaViolationError) Error() string { return e.Err.Error() }

func (e *ProfileSchemaViolationError) Unwrap() error { return e.Err }
