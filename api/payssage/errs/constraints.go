package errs

import "errors"

// One error per named constraint in the schema. internal/app's constraint
// registry maps each constraint to its error, and lib/database pairs it
// with the violation's class. Boot fails when the schema has a constraint
// the registry does not list.
var (
	// intents
	ErrInvalidAmount             = errors.New("amount must be greater than zero")
	ErrInvalidIntentStatus       = errors.New("intent status must be one of: pending, processing, succeeded, cancelled, failed, rejected, refunded")
	ErrInvalidIntentStatusDetail = errors.New("intent status detail must be empty or one of: insufficient_funds, high_risk, invalid_card, card_disabled, expired_card, invalid_security_code, pending_review, other")

	// oauth_states, org_members
	ErrInvalidOAuthFlowType = errors.New("OAuth flow must be one of: collector, seller")
	ErrInvalidOrgRole       = errors.New("organization role must be one of: owner, admin, member")

	// wallets
	ErrNegativeFee = errors.New("fee (bps) must be non-negative")

	// webhooks
	ErrInvalidDeliveryStatus = errors.New("webhook delivery status must be one of: pending, delivered, failed")
)
