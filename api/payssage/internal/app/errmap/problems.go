// Package errmap is Payssage's problem switch (ADR-0019). It sits apart
// from internal/app, which wires the features, so a feature's tests can
// check what status an error answers with without an import cycle.
package errmap

import (
	"errors"
	"net/http"

	"payssage/errs"
	"payssage/internal/problems"

	fp "github.com/MintzyG/problematic"
)

// Problem answers what Payssage's problems can't name for themselves
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
		provider     *errs.ProviderError
		revoked      *errs.SellerRevokedError
		missingParam *errs.MissingParamError
	)
	switch {
	// providers and intents
	case errors.As(err, &provider):
		return providerProblem(provider)
	case errors.As(err, &revoked) && revoked.Op == "checkout":
		return problems.SellerRevokedCheckout(revoked.Error(), revoked.SellerID)
	case errors.As(err, &revoked):
		return problems.SellerRevoked(revoked.Error(), revoked.SellerID)

	// incoming webhooks
	case errors.As(err, &missingParam):
		return problems.Validation(missingParam.Error())
	}
	return nil
}

// providerProblem answers a provider's own error by what the provider said:
// its auth failures are ours to fix (401), a provider outage is a bad
// gateway, and its 4xx answers pass through.
func providerProblem(e *errs.ProviderError) *fp.Problem {
	detail := e.Error()
	switch {
	case e.Status == http.StatusUnauthorized || e.Status == http.StatusForbidden:
		return problems.ProviderUnauthorized(detail, e.Provider, e.Reason)
	case e.Status == http.StatusNotFound:
		return problems.ProviderNotFound(detail, e.Provider, e.Reason)
	case e.Status == http.StatusConflict:
		return problems.ProviderConflict(detail, e.Provider, e.Reason)
	case e.Status == http.StatusTooManyRequests:
		return problems.ProviderRateLimited(detail, e.Provider, e.Reason)
	case e.Status >= http.StatusInternalServerError:
		return problems.ProviderUnavailable(detail, e.Provider, e.Reason)
	default:
		return problems.ProviderRejected(detail, e.Provider, e.Reason)
	}
}
