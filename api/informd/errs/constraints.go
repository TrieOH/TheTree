package errs

import "errors"

// forms
var (
	ErrInvalidFormStatus      = errors.New("status must be one of: draft, open, closed, archived")
	ErrInvalidFormStatusState = errors.New("opened_at, closed_at or archived_at must be set when status is open, closed or archived")
)

// fields
var (
	ErrInvalidFieldType      = errors.New("field type must be one of: string, email, int, float, bool, date, time, datetime, select, file, phone, url")
	ErrInvalidFieldKeyFormat = errors.New("field key must start with a letter or underscore and contain only lowercase letters, digits and underscores")
	ErrInvalidSelectBehavior = errors.New("select fields must have a behaviour of checkbox or radio")
	ErrEmptySelectOptions    = errors.New("select fields must have a non-empty options array")
	ErrInvalidSelectType     = errors.New("select_type must be one of: email, int, float, date, time, datetime, phone, url")
)

// namespaces
var (
	ErrInvalidNamespaceMemberRole = errors.New("invalid namespace member role")
)

// responses, responders and invites
var (
	ErrInviteBothNamed = errors.New("an invite cannot have both responder_id and email set")
)
