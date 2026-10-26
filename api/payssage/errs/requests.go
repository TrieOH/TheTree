package errs

import "errors"

// Request shapes the spec cannot express.
var (
	ErrInvalidRequest              = errors.New("invalid request")
	ErrInvalidProviderData         = errors.New("invalid provider_data")
	ErrInvalidCheckoutProviderData = errors.New("invalid checkout_provider_data")
	ErrInvalidMetadata             = errors.New("invalid metadata")
	ErrTestModeOnly                = errors.New("test mode only")
)

// Ownership outside the role checks (lib/authz.ErrInsufficientRole covers
// those): a personal resource reached by someone other than its owner.
var ErrNotResourceOwner = errors.New("insufficient permissions")

// Organization membership.
var (
	ErrAddSelfToOrg      = errors.New("users can't add themselves to organizations")
	ErrAddOwnerToOrg     = errors.New("owners can't be added to organizations they own")
	ErrAlreadyOrgMember  = errors.New("user is already a member of the organization")
	ErrRemoveSelfFromOrg = errors.New("cannot remove yourself from the organization")
	ErrRemoveOrgOwner    = errors.New("cannot remove the owner of the organization")
	ErrNotOrgMember      = errors.New("user is not a member of the organization")
)

// OAuth connect, callback and revoke.
var (
	ErrInvalidFlow           = errors.New("invalid flow")
	ErrInvalidProvider       = errors.New("invalid provider")
	ErrProviderMismatch      = errors.New("provider mismatch")
	ErrSellerFlowNeedsWallet = errors.New("wallet_id is required for seller flow")
)
