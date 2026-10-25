// Package fielderr attaches the offending request fields to a named error
// (ADR-0016). The error still reads and matches as the named error; the
// HTTP writer (lib/httpserver) lists the fields in the response.
package fielderr

// Field is one offending request field and what is wrong with it.
type Field struct {
	Field   string
	Message string
}

// Error is a named error with the fields that caused it.
type Error struct {
	Err    error
	Fields []Field
}

func (e *Error) Error() string { return e.Err.Error() }

func (e *Error) Unwrap() error { return e.Err }

// Wrap attaches fields to err.
func Wrap(err error, fields ...Field) error {
	return &Error{Err: err, Fields: fields}
}

// On is Wrap for the common single field.
func On(err error, field, message string) error {
	return &Error{Err: err, Fields: []Field{{Field: field, Message: message}}}
}
