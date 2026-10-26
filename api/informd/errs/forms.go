package errs

import "errors"

// Form lifecycle.
var (
	ErrOpenNotDraft            = errors.New("cannot open a form not on draft")
	ErrCloseNotOpen            = errors.New("cannot close a form not on open")
	ErrRedraftNotOpen          = errors.New("cannot redraft a form not on open")
	ErrRedraftWithResponses    = errors.New("cannot redraft a form with responses")
	ErrArchiveNotClosed        = errors.New("cannot archive a form not on closed")
	ErrFormNotAnswerable       = errors.New("form is not open for answers")
	ErrFormNotOpenForResponses = errors.New("form is not open for responses")
)

// Form contents.
var (
	ErrStepsAcrossForms     = errors.New("all steps must belong to the same form")
	ErrFieldsAcrossSteps    = errors.New("all fields must belong to the same step")
	ErrInvalidSelectOptions = errors.New("invalid options")
)
