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
// (ADR-0019): the typed errors, which carry the problem's members, the
// named errors answered with a shared problem, the not-found errors and the
// check constraints. A problem that names its sentinel (x-problem's error)
// is answered by problems.Mapper itself, which calls Problem first and
// tries the shared cases last; Problem is nil for an error it doesn't
// know. Named errors don't wrap named errors, so case order never breaks a
// tie. The detail is the named error's message, never the chain's.
func Problem(err error) *fp.Problem {
	var (
		unsupported  *errs.UnsupportedProviderError
		provider     *errs.ProviderError
		intentStatus *errs.IntentStatusError
		revoked      *errs.SellerRevokedError
		mismatch     *errs.SellerProviderMismatchError
		missingParam *errs.MissingParamError
		unhandled    *errs.UnhandledWebhookEventError
		noIntent     *errs.WebhookIntentNotFoundError
		noToken      *errs.WebhookSellerTokenMissingError
	)
	switch {
	// providers and intents
	case errors.As(err, &unsupported):
		return problems.UnsupportedProvider(unsupported.Error(), unsupported.Provider)
	case errors.As(err, &provider):
		return providerProblem(provider)
	case errors.As(err, &intentStatus):
		return problems.IntentStatus(intentStatus.Error(), intentStatus.Action, intentStatus.Status)
	case errors.As(err, &revoked) && revoked.Op == "checkout":
		return problems.SellerRevokedCheckout(revoked.Error(), revoked.SellerID.String())
	case errors.As(err, &revoked):
		return problems.SellerRevoked(revoked.Error(), revoked.SellerID.String())
	case errors.As(err, &mismatch):
		return problems.SellerProviderMismatch(mismatch.Error(), mismatch.SellerID.String(), mismatch.Provider)

	// incoming webhooks
	case errors.As(err, &missingParam):
		return problems.Validation(missingParam.Error())
	case errors.As(err, &unhandled):
		return problems.UnhandledWebhookEvent(unhandled.Error(), unhandled.Type)
	case errors.As(err, &noIntent):
		return problems.WebhookIntentNotFound(noIntent.Error(), noIntent.PaymentID)
	case errors.As(err, &noToken):
		return problems.WebhookSellerTokenMissing(noToken.Error(), noToken.SellerID.String(), noToken.IntentID.String())

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
	errs.ErrOAuthStateNotFound,
	errs.ErrWebhookEventNotFound,
	errs.ErrWalletNotFound,
	errs.ErrWebhookEndpointNotFound,
	errs.ErrIntentNotFound,
	errs.ErrOrganizationNotFound,
	errs.ErrCollectorNotFound,
	errs.ErrWebhookDeliveryNotFound,
	errs.ErrSellerNotFound,
}

// invalid are the check constraints, answered as the shared validation
// problem.
var invalid = []error{
	errs.ErrInvalidAmount,
	errs.ErrInvalidIntentStatus,
	errs.ErrInvalidIntentStatusDetail,
	errs.ErrInvalidOAuthFlowType,
	errs.ErrInvalidOrgRole,
	errs.ErrNegativeFee,
	errs.ErrInvalidDeliveryStatus,
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
