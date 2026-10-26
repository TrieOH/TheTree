package app

import (
	"IdentityX/errs"
	"lib/database"
)

// Constraints maps every named constraint in the schema to the named error
// repos surface on violation (ADR-0016); the postgres module takes it at
// boot, and boot fails when the schema has a constraint missing here.
func Constraints() database.ConstraintErrors {
	return database.ConstraintErrors{
		// actors
		"chk_actors_type":                      errs.ErrInvalidActorType,
		"chk_actors_auth_method":               errs.ErrInvalidAuthMethod,
		"chk_actors_email_required_for_humans": errs.ErrHumanEmailRequired,
		"uniq_email_per_scope_per_method":      errs.ErrEmailTakenForScope,

		// actor_external_identities
		"chk_actor_external_identities_provider": errs.ErrInvalidExternalIDP,
		"uniq_external_identity_per_actor":       errs.ErrExternalIdentityTaken,

		// organizations
		"uniq_organizations_slug": errs.ErrOrgSlugTaken,

		// org_members
		"chk_org_members_role": errs.ErrInvalidOrgRole,

		// projects
		"uniq_projects_slug":    errs.ErrProjectSlugTaken,
		"uniq_verified_domain":  errs.ErrDomainTaken,
		"chk_brand_slug_format": errs.ErrInvalidBrandSlugFormat,
		"chk_brand_slug_length": errs.ErrInvalidBrandSlugLength,

		// project_members
		"chk_project_members_role": errs.ErrInvalidProjectRole,

		// project_oauth_providers
		"chk_project_oauth_providers_provider": errs.ErrInvalidOAuthProvider,
		"uniq_project_oauth_provider":          errs.ErrOAuthProviderConfigured,

		// oauth_login_states
		"chk_oauth_login_states_provider": errs.ErrInvalidOAuthProvider,
		"uniq_oauth_login_state":          errs.ErrDuplicateOAuthState,

		// platform_roles
		"chk_platform_roles_role": errs.ErrInvalidPlatformRole,

		// api_keys
		"uniq_idx_api_keys_display_prefix": errs.ErrAPIKeyPrefixTaken,

		// capabilities
		"uniq_capability_per_scope": errs.ErrCapabilityExists,

		// profile schemas
		"uniq_project_profile_schema_project_id": errs.ErrProjectProfileSchemaExists,
		"uniq_profile_schema_version_per_scope":  errs.ErrProfileSchemaVersionExists,
		"uniq_profile_schema_version_platform":   errs.ErrPlatformSchemaVersionExists,

		// actor profiles
		"uniq_actor_profiles_handle": errs.ErrHandleTaken,

		// crypto_keys
		"chk_crypto_keys_type":   errs.ErrInvalidCryptoKeyType,
		"chk_crypto_keys_status": errs.ErrInvalidCryptoKeyStat,

		// email_templates
		"chk_email_templates_kind":        errs.ErrInvalidTemplateKindColumn,
		"uniq_email_template_per_project": errs.ErrEmailTemplateExists,

		// terms_of_service
		"chk_terms_of_service_content":      errs.ErrTosContentEmpty,
		"chk_terms_of_service_version":      errs.ErrInvalidTosVersion,
		"uniq_terms_of_service_per_project": errs.ErrTosAlreadyExists,

		// tos_acceptances
		"chk_tos_acceptances_source":            errs.ErrInvalidAcceptanceSource,
		"uniq_tos_acceptance_per_actor_version": errs.ErrTosAlreadyAccepted,

		// action_tokens
		"chk_action_tokens_purpose": errs.ErrInvalidActionTokenPurpose,

		// blacklist_entries
		"chk_blacklist_entries_type":         errs.ErrInvalidBlacklistType,
		"uniq_blacklist_target_type_project": errs.ErrAlreadyBlacklisted,
	}
}
