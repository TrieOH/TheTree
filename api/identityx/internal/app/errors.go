package app

import (
	"net/http"

	"IdentityX/errs"
	"lib/httpserver"
)

// Errors is IdentityX's one table from named errors to HTTP (ADR-0016).
// The most specific error in a chain wins, and an error missing here is a
// generic 500. The shared lib errors (missing rows, constraint classes,
// insufficient role) are mapped by lib/httpserver after these.
func Errors() httpserver.ErrorMap {
	const (
		badRequest   = http.StatusBadRequest
		unauthorized = http.StatusUnauthorized
		forbidden    = http.StatusForbidden
		notFound     = http.StatusNotFound
		conflict     = http.StatusConflict
		unavailable  = http.StatusServiceUnavailable
	)
	is := httpserver.Is
	invalid := func(target error) httpserver.ErrorMapping { return is(target, badRequest).Validation() }

	return httpserver.ErrorMap{
		// identity and tokens
		is(errs.ErrNoIdentity, unauthorized),
		is(errs.ErrPlatformOnly, unauthorized),
		is(errs.ErrInvalidToken, unauthorized),
		is(errs.ErrInvalidAccessToken, unauthorized),
		is(errs.ErrMissingKID, unauthorized),
		is(errs.ErrInvalidKID, unauthorized),
		is(errs.ErrOutdatedToken, unauthorized),
		is(errs.ErrSigningKeyRevoked, unauthorized),
		is(errs.ErrSigningKeyRetired, unauthorized),
		is(errs.ErrTokenRevoked, unauthorized),

		// api keys
		is(errs.ErrInvalidAPIKey, forbidden),
		is(errs.ErrAPIKeyExpired, forbidden),
		is(errs.ErrInvalidCapabilities, badRequest),

		// login, registration and setup
		is(errs.ErrInvalidCredentials, unauthorized),
		is(errs.ErrPasswordOnlyForHumans, badRequest),
		invalid(errs.ErrTermsNotAccepted),
		is(errs.ErrSetupComplete, conflict),
		is(errs.ErrSetupRequired, unavailable),
		is(errs.ErrSetupCheckFailed, unavailable),

		// action tokens
		is(errs.ErrActionTokenInvalid, badRequest),
		is(errs.ErrActionTokenExpired, badRequest),
		is(errs.ErrActionTokenUsed, badRequest),

		// access
		is(errs.ErrActorOutsideProject, forbidden),
		is(errs.ErrNotProfileOwner, forbidden),
		is(errs.ErrActorNotInProject, forbidden),

		// membership
		is(errs.ErrAddSelfToProject, badRequest),
		is(errs.ErrAddOwnerToProject, badRequest),
		is(errs.ErrAlreadyProjectMember, badRequest),
		is(errs.ErrRemoveSelfFromProject, badRequest),
		is(errs.ErrRemoveProjectOwner, badRequest),
		is(errs.ErrNotProjectMember, badRequest),
		is(errs.ErrOrgOwnerAlreadyInProject, badRequest),
		is(errs.ErrOrgMemberAlreadyInProject, badRequest),
		is(errs.ErrRemoveOrgOwnerFromProject, badRequest),
		is(errs.ErrRemoveOrgMemberFromProject, badRequest),
		is(errs.ErrAddSelfToOrg, badRequest),
		is(errs.ErrAddOwnerToOrg, badRequest),
		is(errs.ErrAlreadyOrgMember, badRequest),
		is(errs.ErrRemoveSelfFromOrg, badRequest),
		is(errs.ErrRemoveOrgOwner, badRequest),
		is(errs.ErrNotOrgMember, badRequest),

		// oauth
		is(errs.ErrMissingOAuthCode, badRequest),
		is(errs.ErrMissingOAuthState, badRequest),
		is(errs.ErrOAuthStateNotFound, badRequest),
		is(errs.ErrOAuthStateMismatch, badRequest),
		is(errs.ErrOAuthStateExpired, badRequest),
		is(errs.ErrOAuthCodeExchange, unauthorized),
		is(errs.ErrOAuthEmptyToken, unauthorized),
		is(errs.ErrOAuthUserinfo, unauthorized),
		is(errs.ErrOAuthIncompleteUserinfo, unauthorized),
		is(errs.ErrOAuthGitHubEmail, unauthorized),
		httpserver.As[*errs.UnsupportedProviderError](badRequest),
		httpserver.As[*errs.ProviderNotConfiguredError](badRequest),
		httpserver.As[*errs.ProviderDisabledError](forbidden),
		httpserver.As[*errs.ProviderRemovedError](badRequest),

		// email templates, terms of service, profiles
		invalid(errs.ErrInvalidTemplateKind),
		invalid(errs.ErrTemplateMissingActionURL),
		httpserver.As[*errs.TemplateRenderError](badRequest).Validation(),
		is(errs.ErrTosExists, conflict),
		invalid(errs.ErrTosEmpty),
		httpserver.As[*errs.ProfileSchemaViolationError](badRequest).Validation(),

		// not found
		is(errs.ErrActorNotFound, notFound),
		is(errs.ErrProjectNotFound, notFound),
		is(errs.ErrProfileNotFound, notFound),
		is(errs.ErrExternalIdentityNotFound, notFound),
		is(errs.ErrPlatformRoleNotFound, notFound),
		is(errs.ErrAPIKeyNotFound, notFound),
		is(errs.ErrCryptoKeyNotFound, notFound),
		is(errs.ErrCapabilityNotFound, notFound),
		is(errs.ErrActionTokenNotFound, notFound),
		is(errs.ErrProfileSchemaNotFound, notFound),
		is(errs.ErrBlacklistEntryNotFound, notFound),
		is(errs.ErrOrganizationNotFound, notFound),
		is(errs.ErrOAuthProviderNotFound, notFound),
		is(errs.ErrEmailTemplateNotFound, notFound),
		is(errs.ErrTosNotFound, notFound),

		// constraints: unique → 409, check → 400 validation
		invalid(errs.ErrInvalidActorType),
		invalid(errs.ErrInvalidAuthMethod),
		invalid(errs.ErrHumanEmailRequired),
		is(errs.ErrEmailTakenForScope, conflict),
		invalid(errs.ErrInvalidExternalIDP),
		is(errs.ErrExternalIdentityTaken, conflict),
		is(errs.ErrOrgSlugTaken, conflict),
		invalid(errs.ErrInvalidOrgRole),
		is(errs.ErrProjectSlugTaken, conflict),
		is(errs.ErrDomainTaken, conflict),
		invalid(errs.ErrInvalidBrandSlugFormat),
		invalid(errs.ErrInvalidBrandSlugLength),
		invalid(errs.ErrInvalidProjectRole),
		invalid(errs.ErrInvalidPlatformRole),
		invalid(errs.ErrInvalidOAuthProvider),
		is(errs.ErrOAuthProviderConfigured, conflict),
		is(errs.ErrDuplicateOAuthState, conflict),
		is(errs.ErrAPIKeyPrefixTaken, conflict),
		is(errs.ErrCapabilityExists, conflict),
		invalid(errs.ErrInvalidCryptoKeyType),
		invalid(errs.ErrInvalidCryptoKeyStat),
		is(errs.ErrProjectProfileSchemaExists, conflict),
		is(errs.ErrProfileSchemaVersionExists, conflict),
		is(errs.ErrPlatformSchemaVersionExists, conflict),
		is(errs.ErrHandleTaken, conflict),
		invalid(errs.ErrInvalidTemplateKindColumn),
		is(errs.ErrEmailTemplateExists, conflict),
		invalid(errs.ErrTosContentEmpty),
		invalid(errs.ErrInvalidTosVersion),
		is(errs.ErrTosAlreadyExists, conflict),
		invalid(errs.ErrInvalidAcceptanceSource),
		is(errs.ErrTosAlreadyAccepted, conflict),
		invalid(errs.ErrInvalidActionTokenPurpose),
		invalid(errs.ErrInvalidBlacklistType),
		is(errs.ErrAlreadyBlacklisted, conflict),
	}
}
