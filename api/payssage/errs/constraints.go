package errs

import "errors"

// One error per named constraint in the schema. internal/app's constraint
// registry maps each constraint to its error, and lib/database pairs it
// with the violation's class. Boot fails when the schema has a constraint
// the registry does not list.
var (
	// intents
	ErrInvalidAmount             = errors.New("amount must be greater than zero")
	ErrInvalidIntentStatus       = errors.New("invalid intent status")
	ErrInvalidIntentStatusDetail = errors.New("invalid intent status detail")

	// oauth_states, org_members
	ErrInvalidOAuthFlowType = errors.New("invalid oauth flow type")
	ErrInvalidOrgRole       = errors.New("invalid organization member role")

	// wallets
	ErrNegativeFee = errors.New("fee (bps) must be non-negative")

	// webhooks
	ErrInvalidDeliveryStatus = errors.New("invalid webhook delivery status")
)
