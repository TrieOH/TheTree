package app

import (
	"lib/database"
)

func SetupConstraintMessages() {
	database.SetConstraintErrorRegistry(database.ConstraintRegistry{
		// intents
		"chk_intents_amount_cents":  "amount must be greater than zero",
		"chk_intents_status":        "invalid intent status",
		"chk_intents_status_detail": "invalid intent status detail",

		// oauth_states
		"chk_oauth_states_flow": "invalid oauth flow type",

		// org_members
		"chk_org_members_role": "invalid organization member role",

		// wallets
		"chk_wallets_fee_bps":        "fee (bps) must be non-negative",
		"uniq_wallets_org_name":      "a wallet with this name already exists in this organization",
		"uniq_wallets_personal_name": "a wallet with this name already exists",

		// webhook_deliveries
		"chk_webhook_deliveries_status":          "invalid webhook delivery status",
		"uniq_webhook_deliveries_event_endpoint": "a delivery for this event and endpoint already exists",

		// webhook_events.sql
		"uniq_webhook_events_external_id_type": "a webhook event for this external id and type already exists for this provider",

		// organizations
		"uniq_organizations_slug": "an organization with this slug already exists",

		// provider_credentials
		"uniq_provider_credentials_active": "credentials for this provider are already connected to this wallet",

		// collectors
		"uniq_collectors_org_active":      "this collector is already connected to this organization",
		"uniq_collectors_personal_active": "this collector is already connected to your account",

		// sellers
		"uniq_sellers_active": "this seller is already connected to this wallet",
	})
}
