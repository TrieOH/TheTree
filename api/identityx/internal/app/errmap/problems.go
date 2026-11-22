// Package errmap is IdentityX's problem switch (ADR-0019). It sits apart
// from internal/app, which wires the features, so a feature's tests can
// check what status an error answers with without an import cycle.
package errmap

import (
	"errors"

	"IdentityX/errs"
	"IdentityX/internal/problems"

	fp "github.com/MintzyG/problematic"
)

// Problem answers what IdentityX's problems can't name for themselves
// (ADR-0019): the typed errors, which carry the problem's members, and the
// named errors answered with a shared problem outside x-shared-errors. A
// problem that names its sentinel (x-problem's error) and an error listed
// in x-shared-errors are answered by problems.Mapper itself, which calls
// Problem first and tries the shared cases last; Problem is nil for an
// error it doesn't know. Named errors don't wrap named errors, so case
// order never breaks a tie. The detail is the named error's message, never
// the chain's.
func Problem(err error) *fp.Problem {
	var (
		notConfigured *errs.ProviderNotConfiguredError
		render        *errs.TemplateRenderError
		violation     *errs.ProfileSchemaViolationError
	)
	switch {
	// identity and tokens
	case errors.Is(err, errs.ErrNoIdentity):
		return problems.Unauthenticated(errs.ErrNoIdentity.Error())
	case errors.Is(err, errs.ErrInvalidToken):
		return problems.InvalidToken(errs.ErrInvalidToken.Error())
	case errors.Is(err, errs.ErrTokenExpired):
		return problems.TokenExpired(errs.ErrTokenExpired.Error())

	// api keys
	case errors.Is(err, errs.ErrInvalidAPIKey):
		return problems.InvalidAPIKey(errs.ErrInvalidAPIKey.Error())

	// oauth
	case errors.As(err, &notConfigured):
		return problems.ProviderNotConfigured(notConfigured.Error(), notConfigured.Provider)

	// email templates, terms of service, profiles
	case errors.As(err, &render):
		return problems.TemplateRender(render.Error())
	case errors.As(err, &violation):
		return problems.ProfileSchemaViolation(violation.Error())

	}
	return nil
}
