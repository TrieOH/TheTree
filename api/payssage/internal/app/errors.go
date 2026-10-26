package app

import (
	"net/http"

	"lib/httpserver"
	"payssage/errs"
)

// Errors is Payssage's one table from named errors to HTTP (ADR-0016). The
// most specific error in a chain wins, and an error missing here is a
// generic 500. The shared lib errors (missing rows, constraint classes,
// batches, insufficient role) are mapped by lib/httpserver after these.
func Errors() httpserver.ErrorMap {
	const (
		badRequest   = http.StatusBadRequest
		unauthorized = http.StatusUnauthorized
		forbidden    = http.StatusForbidden
		notFound     = http.StatusNotFound
		conflict     = http.StatusConflict
		badGateway   = http.StatusBadGateway
		unavailable  = http.StatusServiceUnavailable
	)
	is := httpserver.Is
	invalid := func(target error) httpserver.ErrorMapping { return is(target, badRequest).Validation() }

	return httpserver.ErrorMap{
		// request shapes
		is(errs.ErrInvalidRequest, badRequest),
		is(errs.ErrInvalidProviderData, badRequest),
		is(errs.ErrInvalidCheckoutProviderData, badRequest),
		is(errs.ErrInvalidMetadata, badRequest),
		is(errs.ErrTestModeOnly, unavailable),
		is(errs.ErrNotResourceOwner, forbidden),

		// organization membership
		is(errs.ErrAddSelfToOrg, badRequest),
		is(errs.ErrAddOwnerToOrg, badRequest),
		is(errs.ErrAlreadyOrgMember, badRequest),
		is(errs.ErrRemoveSelfFromOrg, badRequest),
		is(errs.ErrRemoveOrgOwner, badRequest),
		is(errs.ErrNotOrgMember, badRequest),

		// oauth
		is(errs.ErrInvalidFlow, badRequest),
		is(errs.ErrInvalidProvider, badRequest),
		is(errs.ErrProviderMismatch, badRequest),
		is(errs.ErrSellerFlowNeedsWallet, badRequest),

		// providers and intents
		is(errs.ErrProviderNotImplemented, badRequest),
		is(errs.ErrProviderRequestFailed, badGateway),
		httpserver.As[*errs.UnsupportedProviderError](badRequest),
		httpserver.AsStatus(providerStatus),
		is(errs.ErrSellerNotInWallet, badRequest),
		is(errs.ErrNoTransactionToRefund, conflict),
		is(errs.ErrNoTransactionToCancel, conflict),
		httpserver.As[*errs.IntentStatusError](badRequest),
		httpserver.AsStatus(func(e *errs.SellerRevokedError) int {
			if e.Op == "checkout" {
				return forbidden
			}
			return conflict
		}),
		httpserver.As[*errs.SellerProviderMismatchError](conflict),

		// incoming webhooks
		is(errs.ErrUnknownWebhookProvider, notFound),
		is(errs.ErrWebhookNotImplemented, badRequest),
		is(errs.ErrWebhookPayloadInvalid, badRequest),
		is(errs.ErrWebhookMissingDataID, badRequest),
		is(errs.ErrWebhookSignatureMissing, unauthorized),
		is(errs.ErrWebhookSignatureMismatch, unauthorized),
		httpserver.As[*errs.MissingParamError](badRequest).Validation(),
		httpserver.As[*errs.UnhandledWebhookEventError](notFound),
		httpserver.As[*errs.WebhookIntentNotFoundError](notFound),
		httpserver.As[*errs.WebhookSellerTokenMissingError](conflict),

		// not found
		is(errs.ErrOAuthStateNotFound, notFound),
		is(errs.ErrWebhookEventNotFound, notFound),
		is(errs.ErrWalletNotFound, notFound),
		is(errs.ErrWebhookEndpointNotFound, notFound),
		is(errs.ErrIntentNotFound, notFound),
		is(errs.ErrOrganizationNotFound, notFound),
		is(errs.ErrCollectorNotFound, notFound),
		is(errs.ErrWebhookDeliveryNotFound, notFound),
		is(errs.ErrSellerNotFound, notFound),

		// constraints: unique → 409, check → 400 validation
		invalid(errs.ErrInvalidAmount),
		invalid(errs.ErrInvalidIntentStatus),
		invalid(errs.ErrInvalidIntentStatusDetail),
		invalid(errs.ErrInvalidOAuthFlowType),
		invalid(errs.ErrInvalidOrgRole),
		invalid(errs.ErrNegativeFee),
		is(errs.ErrOrgWalletNameTaken, conflict),
		is(errs.ErrPersonalWalletTaken, conflict),
		invalid(errs.ErrInvalidDeliveryStatus),
		is(errs.ErrDeliveryExists, conflict),
		is(errs.ErrWebhookEventExists, conflict),
		is(errs.ErrOrgSlugTaken, conflict),
		is(errs.ErrCredentialsConnected, conflict),
		is(errs.ErrOrgCollectorConnected, conflict),
		is(errs.ErrPersonalCollectorConnected, conflict),
		is(errs.ErrSellerConnected, conflict),
	}
}

// providerStatus answers a provider's own error by what the provider said:
// its auth failures are ours to fix (401), a provider outage is a bad
// gateway, and its 4xx answers pass through.
func providerStatus(e *errs.ProviderError) int {
	switch {
	case e.Status == http.StatusUnauthorized || e.Status == http.StatusForbidden:
		return http.StatusUnauthorized
	case e.Status == http.StatusNotFound:
		return http.StatusNotFound
	case e.Status == http.StatusConflict:
		return http.StatusConflict
	case e.Status == http.StatusTooManyRequests:
		return http.StatusTooManyRequests
	case e.Status >= http.StatusInternalServerError:
		return http.StatusBadGateway
	default:
		return http.StatusBadRequest
	}
}
