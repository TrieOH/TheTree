package errs

import "errors"

// Email templates.
var (
	ErrInvalidTemplateKind      = errors.New("email template kind must be one of: verify, reset, tos")
	ErrTemplateMissingActionURL = errors.New("email template body must include the {{.ActionURL}} variable at least once")
)

// TemplateRenderError is an email template override that does not render.
type TemplateRenderError struct{ Err error }

func (e *TemplateRenderError) Error() string {
	return "email template does not render: " + e.Err.Error()
}

func (e *TemplateRenderError) Unwrap() error { return e.Err }

// Terms of service.
var (
	ErrTosExists = errors.New("project already has terms of service; update them instead")
	ErrTosEmpty  = errors.New("terms of service content must not be empty")
)

// ProfileSchemaViolationError is a profile its scope's schema rejects; the
// message is the validator's.
type ProfileSchemaViolationError struct{ Err error }

func (e *ProfileSchemaViolationError) Error() string { return e.Err.Error() }

func (e *ProfileSchemaViolationError) Unwrap() error { return e.Err }
