package errs

import "errors"

// Each resource's not-found error. Repos return them joined with
// lib/database.ErrNotFound, so a caller can match the resource or any
// missing row.
var (
	ErrFormNotFound      = errors.New("form not found")
	ErrNamespaceNotFound = errors.New("namespace not found")
	ErrStepNotFound      = errors.New("step not found")
	ErrFieldNotFound     = errors.New("field not found")
	ErrAnswerNotFound    = errors.New("answer not found")
	ErrResponseNotFound  = errors.New("response not found")
	ErrResponderNotFound = errors.New("responder not found")
)
