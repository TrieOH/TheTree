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
	ErrNegativeFee         = errors.New("fee (bps) must be non-negative")
	ErrOrgWalletNameTaken  = errors.New("a wallet with this name already exists in this organization")
	ErrPersonalWalletTaken = errors.New("a wallet with this name already exists")

	// webhooks
	ErrInvalidDeliveryStatus = errors.New("invalid webhook delivery status")
	ErrDeliveryExists        = errors.New("a delivery for this event and endpoint already exists")
	ErrWebhookEventExists    = errors.New("a webhook event for this external id and type already exists for this provider")

	// organizations, credentials, collectors, sellers
	ErrOrgSlugTaken               = errors.New("an organization with this slug already exists")
	ErrCredentialsConnected       = errors.New("credentials for this provider are already connected to this wallet")
	ErrOrgCollectorConnected      = errors.New("this collector is already connected to this organization")
	ErrPersonalCollectorConnected = errors.New("this collector is already connected to your account")
	ErrSellerConnected            = errors.New("this seller is already connected to this wallet")
)
