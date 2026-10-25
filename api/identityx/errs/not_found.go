package errs

import "errors"

// Each resource's not-found error. Repos return them joined with
// lib/database.ErrNotFound, so a caller can match the resource or any
// missing row.
var (
	ErrActorNotFound            = errors.New("actor not found")
	ErrProjectNotFound          = errors.New("project not found")
	ErrProfileNotFound          = errors.New("profile not found")
	ErrExternalIdentityNotFound = errors.New("external actor identity not found")
	ErrPlatformRoleNotFound     = errors.New("platform role not found")
	ErrAPIKeyNotFound           = errors.New("api key not found")
	ErrCryptoKeyNotFound        = errors.New("crypto key not found")
	ErrCapabilityNotFound       = errors.New("capability not found")
	ErrActionTokenNotFound      = errors.New("action token not found")
	ErrProfileSchemaNotFound    = errors.New("profile schema not found")
	ErrBlacklistEntryNotFound   = errors.New("blacklist entry not found")
	ErrOrganizationNotFound     = errors.New("organization not found")
	ErrOAuthProviderNotFound    = errors.New("project oauth provider not found")
	ErrEmailTemplateNotFound    = errors.New("email template not found")
	ErrTosNotFound              = errors.New("terms of service not found")
)
