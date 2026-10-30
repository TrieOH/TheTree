package app

import (
	"errors"

	"IdentityX/errs"
	"IdentityX/internal/problems"

	fp "github.com/MintzyG/problematic"
)

// Problem answers IdentityX's named errors with their declared problems
// (ADR-0019), and nil for an error it doesn't know; problems.Mapper then
// tries the shared cases. Named errors don't wrap named errors, so case
// order never breaks a tie. The detail is the named error's message, never
// the chain's.
func Problem(err error) *fp.Problem {
	var (
		unsupported   *errs.UnsupportedProviderError
		notConfigured *errs.ProviderNotConfiguredError
		disabled      *errs.ProviderDisabledError
		removed       *errs.ProviderRemovedError
		render        *errs.TemplateRenderError
		violation     *errs.ProfileSchemaViolationError
	)
	switch {
	// identity and tokens
	case errors.Is(err, errs.ErrNoIdentity):
		return problems.Unauthenticated(errs.ErrNoIdentity.Error())
	case errors.Is(err, errs.ErrPlatformOnly):
		return problems.PlatformOnly(errs.ErrPlatformOnly.Error())
	case errors.Is(err, errs.ErrInvalidToken):
		return problems.InvalidToken(errs.ErrInvalidToken.Error())
	case errors.Is(err, errs.ErrInvalidAccessToken):
		return problems.InvalidAccessToken(errs.ErrInvalidAccessToken.Error())
	case errors.Is(err, errs.ErrMissingKID):
		return problems.MissingKID(errs.ErrMissingKID.Error())
	case errors.Is(err, errs.ErrInvalidKID):
		return problems.InvalidKID(errs.ErrInvalidKID.Error())
	case errors.Is(err, errs.ErrOutdatedToken):
		return problems.OutdatedToken(errs.ErrOutdatedToken.Error())
	case errors.Is(err, errs.ErrSigningKeyRevoked):
		return problems.SigningKeyRevoked(errs.ErrSigningKeyRevoked.Error())
	case errors.Is(err, errs.ErrSigningKeyRetired):
		return problems.SigningKeyRetired(errs.ErrSigningKeyRetired.Error())
	case errors.Is(err, errs.ErrTokenRevoked):
		return problems.TokenRevoked(errs.ErrTokenRevoked.Error())

	// api keys
	case errors.Is(err, errs.ErrInvalidAPIKey):
		return problems.InvalidAPIKey(errs.ErrInvalidAPIKey.Error())
	case errors.Is(err, errs.ErrAPIKeyExpired):
		return problems.APIKeyExpired(errs.ErrAPIKeyExpired.Error())
	case errors.Is(err, errs.ErrInvalidCapabilities):
		return problems.InvalidCapabilities(errs.ErrInvalidCapabilities.Error())

	// login, registration and setup
	case errors.Is(err, errs.ErrInvalidCredentials):
		return problems.InvalidCredentials(errs.ErrInvalidCredentials.Error())
	case errors.Is(err, errs.ErrPasswordOnlyForHumans):
		return problems.PasswordOnlyForHumans(errs.ErrPasswordOnlyForHumans.Error())
	case errors.Is(err, errs.ErrSetupComplete):
		return problems.SetupComplete(errs.ErrSetupComplete.Error())
	case errors.Is(err, errs.ErrSetupRequired):
		return problems.SetupRequired(errs.ErrSetupRequired.Error())
	case errors.Is(err, errs.ErrSetupCheckFailed):
		return problems.SetupCheckFailed(errs.ErrSetupCheckFailed.Error())

	// action tokens
	case errors.Is(err, errs.ErrActionTokenInvalid):
		return problems.ActionTokenInvalid(errs.ErrActionTokenInvalid.Error())
	case errors.Is(err, errs.ErrActionTokenExpired):
		return problems.ActionTokenExpired(errs.ErrActionTokenExpired.Error())
	case errors.Is(err, errs.ErrActionTokenUsed):
		return problems.ActionTokenUsed(errs.ErrActionTokenUsed.Error())

	// access
	case errors.Is(err, errs.ErrActorOutsideProject):
		return problems.ActorOutsideProject(errs.ErrActorOutsideProject.Error())
	case errors.Is(err, errs.ErrNotProfileOwner):
		return problems.NotProfileOwner(errs.ErrNotProfileOwner.Error())
	case errors.Is(err, errs.ErrActorNotInProject):
		return problems.ActorNotInProject(errs.ErrActorNotInProject.Error())

	// membership
	case errors.Is(err, errs.ErrAddSelfToProject):
		return problems.AddSelfToProject(errs.ErrAddSelfToProject.Error())
	case errors.Is(err, errs.ErrAddOwnerToProject):
		return problems.AddOwnerToProject(errs.ErrAddOwnerToProject.Error())
	case errors.Is(err, errs.ErrAlreadyProjectMember):
		return problems.AlreadyProjectMember(errs.ErrAlreadyProjectMember.Error())
	case errors.Is(err, errs.ErrRemoveSelfFromProject):
		return problems.RemoveSelfFromProject(errs.ErrRemoveSelfFromProject.Error())
	case errors.Is(err, errs.ErrRemoveProjectOwner):
		return problems.RemoveProjectOwner(errs.ErrRemoveProjectOwner.Error())
	case errors.Is(err, errs.ErrNotProjectMember):
		return problems.NotProjectMember(errs.ErrNotProjectMember.Error())
	case errors.Is(err, errs.ErrOrgOwnerAlreadyInProject):
		return problems.OrgOwnerAlreadyInProject(errs.ErrOrgOwnerAlreadyInProject.Error())
	case errors.Is(err, errs.ErrOrgMemberAlreadyInProject):
		return problems.OrgMemberAlreadyInProject(errs.ErrOrgMemberAlreadyInProject.Error())
	case errors.Is(err, errs.ErrRemoveOrgOwnerFromProject):
		return problems.RemoveOrgOwnerFromProject(errs.ErrRemoveOrgOwnerFromProject.Error())
	case errors.Is(err, errs.ErrRemoveOrgMemberFromProject):
		return problems.RemoveOrgMemberFromProject(errs.ErrRemoveOrgMemberFromProject.Error())
	case errors.Is(err, errs.ErrAddSelfToOrg):
		return problems.AddSelfToOrg(errs.ErrAddSelfToOrg.Error())
	case errors.Is(err, errs.ErrAddOwnerToOrg):
		return problems.AddOwnerToOrg(errs.ErrAddOwnerToOrg.Error())
	case errors.Is(err, errs.ErrAlreadyOrgMember):
		return problems.AlreadyOrgMember(errs.ErrAlreadyOrgMember.Error())
	case errors.Is(err, errs.ErrRemoveSelfFromOrg):
		return problems.RemoveSelfFromOrg(errs.ErrRemoveSelfFromOrg.Error())
	case errors.Is(err, errs.ErrRemoveOrgOwner):
		return problems.RemoveOrgOwner(errs.ErrRemoveOrgOwner.Error())
	case errors.Is(err, errs.ErrNotOrgMember):
		return problems.NotOrgMember(errs.ErrNotOrgMember.Error())

	// oauth
	case errors.Is(err, errs.ErrMissingOAuthCode):
		return problems.MissingOAuthCode(errs.ErrMissingOAuthCode.Error())
	case errors.Is(err, errs.ErrMissingOAuthState):
		return problems.MissingOAuthState(errs.ErrMissingOAuthState.Error())
	case errors.Is(err, errs.ErrOAuthStateNotFound):
		return problems.OAuthStateNotFound(errs.ErrOAuthStateNotFound.Error())
	case errors.Is(err, errs.ErrOAuthStateMismatch):
		return problems.OAuthStateMismatch(errs.ErrOAuthStateMismatch.Error())
	case errors.Is(err, errs.ErrOAuthStateExpired):
		return problems.OAuthStateExpired(errs.ErrOAuthStateExpired.Error())
	case errors.Is(err, errs.ErrOAuthCodeExchange):
		return problems.OAuthCodeExchange(errs.ErrOAuthCodeExchange.Error())
	case errors.Is(err, errs.ErrOAuthEmptyToken):
		return problems.OAuthEmptyToken(errs.ErrOAuthEmptyToken.Error())
	case errors.Is(err, errs.ErrOAuthUserinfo):
		return problems.OAuthUserinfo(errs.ErrOAuthUserinfo.Error())
	case errors.Is(err, errs.ErrOAuthIncompleteUserinfo):
		return problems.OAuthIncompleteUserinfo(errs.ErrOAuthIncompleteUserinfo.Error())
	case errors.Is(err, errs.ErrOAuthGitHubEmail):
		return problems.OAuthGitHubEmail(errs.ErrOAuthGitHubEmail.Error())
	case errors.As(err, &unsupported):
		return problems.UnsupportedProvider(unsupported.Error(), unsupported.Provider)
	case errors.As(err, &notConfigured):
		return problems.ProviderNotConfigured(notConfigured.Error(), notConfigured.Provider)
	case errors.As(err, &disabled):
		return problems.ProviderDisabled(disabled.Error(), disabled.Provider)
	case errors.As(err, &removed):
		return problems.ProviderRemoved(removed.Error(), removed.Provider)

	// email templates, terms of service, profiles
	case errors.As(err, &render):
		return problems.TemplateRender(render.Error())
	case errors.Is(err, errs.ErrTosExists):
		return problems.TosExists(errs.ErrTosExists.Error())
	case errors.As(err, &violation):
		return problems.ProfileSchemaViolation(violation.Error())

	// constraints
	case errors.Is(err, errs.ErrEmailTakenForScope):
		return problems.EmailTakenForScope(errs.ErrEmailTakenForScope.Error())
	case errors.Is(err, errs.ErrExternalIdentityTaken):
		return problems.ExternalIdentityTaken(errs.ErrExternalIdentityTaken.Error())
	case errors.Is(err, errs.ErrOrgSlugTaken):
		return problems.OrgSlugTaken(errs.ErrOrgSlugTaken.Error())
	case errors.Is(err, errs.ErrProjectSlugTaken):
		return problems.ProjectSlugTaken(errs.ErrProjectSlugTaken.Error())
	case errors.Is(err, errs.ErrDomainTaken):
		return problems.DomainTaken(errs.ErrDomainTaken.Error())
	case errors.Is(err, errs.ErrOAuthProviderConfigured):
		return problems.OAuthProviderConfigured(errs.ErrOAuthProviderConfigured.Error())
	case errors.Is(err, errs.ErrDuplicateOAuthState):
		return problems.DuplicateOAuthState(errs.ErrDuplicateOAuthState.Error())
	case errors.Is(err, errs.ErrAPIKeyPrefixTaken):
		return problems.APIKeyPrefixTaken(errs.ErrAPIKeyPrefixTaken.Error())
	case errors.Is(err, errs.ErrCapabilityExists):
		return problems.CapabilityExists(errs.ErrCapabilityExists.Error())
	case errors.Is(err, errs.ErrProjectProfileSchemaExists):
		return problems.ProjectProfileSchemaExists(errs.ErrProjectProfileSchemaExists.Error())
	case errors.Is(err, errs.ErrProfileSchemaVersionExists):
		return problems.ProfileSchemaVersionExists(errs.ErrProfileSchemaVersionExists.Error())
	case errors.Is(err, errs.ErrPlatformSchemaVersionExists):
		return problems.PlatformSchemaVersionExists(errs.ErrPlatformSchemaVersionExists.Error())
	case errors.Is(err, errs.ErrHandleTaken):
		return problems.HandleTaken(errs.ErrHandleTaken.Error())
	case errors.Is(err, errs.ErrEmailTemplateExists):
		return problems.EmailTemplateExists(errs.ErrEmailTemplateExists.Error())
	case errors.Is(err, errs.ErrTosAlreadyExists):
		return problems.TosAlreadyExists(errs.ErrTosAlreadyExists.Error())
	case errors.Is(err, errs.ErrTosAlreadyAccepted):
		return problems.TosAlreadyAccepted(errs.ErrTosAlreadyAccepted.Error())
	case errors.Is(err, errs.ErrAlreadyBlacklisted):
		return problems.AlreadyBlacklisted(errs.ErrAlreadyBlacklisted.Error())
	}
	for _, target := range notFound {
		if errors.Is(err, target) {
			return problems.NotFound(target.Error())
		}
	}
	for _, target := range invalid {
		if errors.Is(err, target) {
			return problems.Validation(target.Error())
		}
	}
	return nil
}

// notFound are the named missing resources, answered as the shared
// not-found problem.
var notFound = []error{
	errs.ErrActorNotFound,
	errs.ErrProjectNotFound,
	errs.ErrProfileNotFound,
	errs.ErrExternalIdentityNotFound,
	errs.ErrPlatformRoleNotFound,
	errs.ErrAPIKeyNotFound,
	errs.ErrCryptoKeyNotFound,
	errs.ErrCapabilityNotFound,
	errs.ErrActionTokenNotFound,
	errs.ErrProfileSchemaNotFound,
	errs.ErrBlacklistEntryNotFound,
	errs.ErrOrganizationNotFound,
	errs.ErrOAuthProviderNotFound,
	errs.ErrEmailTemplateNotFound,
	errs.ErrTosNotFound,
}

// invalid are the check constraints, answered as the shared validation
// problem.
var invalid = []error{
	errs.ErrTermsNotAccepted,
	errs.ErrInvalidTemplateKind,
	errs.ErrTemplateMissingActionURL,
	errs.ErrTosEmpty,
	errs.ErrInvalidActorType,
	errs.ErrInvalidAuthMethod,
	errs.ErrHumanEmailRequired,
	errs.ErrInvalidExternalIDP,
	errs.ErrInvalidOrgRole,
	errs.ErrInvalidBrandSlugFormat,
	errs.ErrInvalidBrandSlugLength,
	errs.ErrInvalidProjectRole,
	errs.ErrInvalidPlatformRole,
	errs.ErrInvalidOAuthProvider,
	errs.ErrInvalidCryptoKeyType,
	errs.ErrInvalidCryptoKeyStat,
	errs.ErrInvalidTemplateKindColumn,
	errs.ErrTosContentEmpty,
	errs.ErrInvalidTosVersion,
	errs.ErrInvalidAcceptanceSource,
	errs.ErrInvalidActionTokenPurpose,
	errs.ErrInvalidBlacklistType,
}
