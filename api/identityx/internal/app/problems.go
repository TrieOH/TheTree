package app

import (
	"errors"

	"IdentityX/errs"
	"IdentityX/internal/problems"

	fp "github.com/MintzyG/problematic"
)

// Problem answers what IdentityX's problems can't name for themselves
// (ADR-0019): the typed errors, which carry the problem's members, the
// named errors answered with a shared problem, the not-found errors and the
// check constraints. A problem that names its sentinel (x-problem's error)
// is answered by problems.Mapper itself, which calls Problem first and
// tries the shared cases last; Problem is nil for an error it doesn't
// know. Named errors don't wrap named errors, so case order never breaks a
// tie. The detail is the named error's message, never the chain's.
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
	case errors.Is(err, errs.ErrInvalidToken):
		return problems.InvalidToken(errs.ErrInvalidToken.Error())

	// api keys
	case errors.Is(err, errs.ErrInvalidAPIKey):
		return problems.InvalidAPIKey(errs.ErrInvalidAPIKey.Error())

	// oauth
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
	case errors.As(err, &violation):
		return problems.ProfileSchemaViolation(violation.Error())

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
