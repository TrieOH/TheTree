package errs

import "errors"

// Each resource's not-found error. Repos return them as
// lib/database.NotFoundError, so a caller can match the resource or any
// missing row.
var (
	ErrEditionNotFound       = errors.New("edition not found")
	ErrBadgeTemplateNotFound = errors.New("badge template not found")
	ErrProgramNotFound       = errors.New("program not found")
	ErrPurchaseNotFound      = errors.New("purchase not found")
	ErrTicketTypeNotFound    = errors.New("ticket type not found")
	ErrRegistrationNotFound  = errors.New("registration not found")
	ErrCertificationNotFound = errors.New("certification not found")
	ErrEventNotFound         = errors.New("event not found")
	ErrProductNotFound       = errors.New("product not found")
	ErrSignatureNotFound     = errors.New("signature not found")
	ErrWSTokenNotFound       = errors.New("ws token not found")
)
