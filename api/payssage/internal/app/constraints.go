package app

import (
	"lib/database"
	"payssage/errs"
)

// Constraints maps every named constraint in the schema to the named error
// repos surface on violation (ADR-0016); the postgres module takes it at
// boot, and boot fails when the schema has a constraint missing here.
func Constraints() database.ConstraintErrors {
	return database.ConstraintErrors{
		// intents
		"chk_intents_amount_cents":  errs.ErrInvalidAmount,
		"chk_intents_status":        errs.ErrInvalidIntentStatus,
		"chk_intents_status_detail": errs.ErrInvalidIntentStatusDetail,

		// oauth_states
		"chk_oauth_states_flow": errs.ErrInvalidOAuthFlowType,

		// org_members
		"chk_org_members_role": errs.ErrInvalidOrgRole,

		// wallets
		"chk_wallets_fee_bps":        errs.ErrNegativeFee,
		"uniq_wallets_org_name":      errs.ErrOrgWalletNameTaken,
		"uniq_wallets_personal_name": errs.ErrPersonalWalletTaken,

		// webhook_deliveries
		"chk_webhook_deliveries_status":          errs.ErrInvalidDeliveryStatus,
		"uniq_webhook_deliveries_event_endpoint": errs.ErrDeliveryExists,

		// webhook_events
		"uniq_webhook_events_external_id_type": errs.ErrWebhookEventExists,

		// organizations
		"uniq_organizations_slug": errs.ErrOrgSlugTaken,

		// provider_credentials
		"uniq_provider_credentials_active": errs.ErrCredentialsConnected,

		// collectors
		"uniq_collectors_org_active":      errs.ErrOrgCollectorConnected,
		"uniq_collectors_personal_active": errs.ErrPersonalCollectorConnected,

		// sellers
		"uniq_sellers_active": errs.ErrSellerConnected,
	}
}
