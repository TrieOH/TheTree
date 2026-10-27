package errs

import "errors"

// One error per named constraint in the schema. internal/app's constraint
// registry maps each constraint to its error, and lib/database joins it
// with the violation's class. Boot fails when the schema has a constraint
// the registry does not list.

// forms
var (
	ErrInvalidFormStatus      = errors.New("status must be one of: draft, open, closed, archived")
	ErrInvalidFormStatusState = errors.New("opened_at, closed_at or archived_at must be set when status is open, closed or archived")
	ErrFormNameTaken          = errors.New("a form with this name already exists in this namespace")
	ErrAPIKeyNameTaken        = errors.New("an API key with this name already exists")
)

// fields
var (
	ErrInvalidFieldType      = errors.New("field type must be one of: string, email, int, float, bool, date, time, datetime, select, file, phone, url")
	ErrInvalidFieldKeyFormat = errors.New("field key must start with a letter or underscore and contain only lowercase letters, digits and underscores")
	ErrFieldKeyTaken         = errors.New("a field with this key already exists in this step")
	ErrInvalidSelectBehavior = errors.New("select fields must have a behaviour of checkbox or radio")
	ErrEmptySelectOptions    = errors.New("select fields must have a non-empty options array")
	ErrInvalidSelectType     = errors.New("select_type must be one of: email, int, float, date, time, datetime, phone, url")
)

// namespaces
var (
	ErrNamespaceNameTaken         = errors.New("a namespace with this name already exists")
	ErrInvalidNamespaceMemberRole = errors.New("invalid namespace member role")
)

// responses, responders and invites
var (
	ErrFieldAlreadyAnswered = errors.New("this field already has an answer in this response")
	ErrResponderEmailTaken  = errors.New("a responder with this email already exists on this system")
	ErrInviteBothNamed      = errors.New("an invite cannot have both responder_id and email set")
	ErrInviteTokenTaken     = errors.New("an invite with this token already exists")
)
