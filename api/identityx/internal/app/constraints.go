package app

import "lib/database"

func SetupConstraintMessages() {
	database.SetConstraintErrorRegistry(database.ConstraintRegistry{
		// actors
		"chk_actors_type":                      "actor type must be one of: human, service, machine",
		"chk_actors_auth_method":               "auth method must be one of: api_key, password, google, github",
		"chk_actors_email_required_for_humans": "email is required for human actors",
		"uniq_email_per_scope_per_method":      "an account with this email already exists for this scope and auth method",

		// actor_external_identities
		"chk_actor_external_identities_provider": "external identity provider must be one of: google, github",
		"uniq_external_identity":                 "this external identity is already linked to another account",

		// organizations
		"uniq_organizations_slug": "an organization with this slug already exists",

		// org_members
		"chk_org_members_role": "organization role must be one of: owner, admin, member",

		// projects
		"uniq_projects_slug":    "a project with this slug already exists",
		"uniq_verified_domain":  "this domain is already verified by another project",
		"chk_brand_slug_format": "brand slug must contain only lowercase letters",
		"chk_brand_slug_length": "brand slug length must be between 3 and 32 characters",

		// project_members
		"chk_project_members_role": "project role must be one of: owner, admin, member",

		// project_oauth_providers
		"chk_project_oauth_providers_provider": "OAuth provider must be one of: google, github",
		"uniq_project_oauth_provider":          "this OAuth provider is already configured for this project",

		// platform_roles
		"chk_platform_roles_role": "platform role must be one of: super_admin, admin, support",

		// api_keys
		"uniq_idx_api_keys_display_prefix": "an API key with this display prefix already exists",

		// capabilities
		"uniq_capability_per_scope": "a capability with this resource and action already exists for this scope",

		// profile schemas
		"uniq_project_profile_schema_project_id": "a profile schema already exists for this project",

		// crypto_keys
		"chk_crypto_keys_type":   "crypto key type must be one of: encryption, signing",
		"chk_crypto_keys_status": "crypto key status must be one of: active, retiring, retired, revoked",

		// blacklist_entries
		"chk_blacklist_entries_type":         "blacklist entry type must be one of: actor, token, api_key, email, ip",
		"uniq_blacklist_target_type_project": "this target is already blacklisted for this scope",
	})
}
