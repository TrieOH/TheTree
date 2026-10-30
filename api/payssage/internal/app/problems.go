package app

import (
	"errors"
	"net/http"

	"payssage/errs"
	"payssage/internal/problems"

	fp "github.com/MintzyG/problematic"
)

// Problem answers Payssage's named errors with their declared problems
// (ADR-0019), and nil for an error it doesn't know; problems.Mapper then
// tries the shared cases. Named errors don't wrap named errors, so case
// order never breaks a tie. The detail is the named error's message, never
// the chain's.
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
	// request shapes
	case errors.Is(err, errs.ErrInvalidRequest):
		return problems.InvalidRequest(errs.ErrInvalidRequest.Error())
	case errors.Is(err, errs.ErrInvalidProviderData):
		return problems.InvalidProviderData(errs.ErrInvalidProviderData.Error())
	case errors.Is(err, errs.ErrInvalidCheckoutProviderData):
		return problems.InvalidCheckoutProviderData(errs.ErrInvalidCheckoutProviderData.Error())
	case errors.Is(err, errs.ErrInvalidMetadata):
		return problems.InvalidMetadata(errs.ErrInvalidMetadata.Error())
	case errors.Is(err, errs.ErrTestModeOnly):
		return problems.TestModeOnly(errs.ErrTestModeOnly.Error())
	case errors.Is(err, errs.ErrNotResourceOwner):
		return problems.NotResourceOwner(errs.ErrNotResourceOwner.Error())

	// organization membership
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
	case errors.Is(err, errs.ErrInvalidFlow):
		return problems.InvalidFlow(errs.ErrInvalidFlow.Error())
	case errors.Is(err, errs.ErrInvalidProvider):
		return problems.InvalidProvider(errs.ErrInvalidProvider.Error())
	case errors.Is(err, errs.ErrProviderMismatch):
		return problems.ProviderMismatch(errs.ErrProviderMismatch.Error())
	case errors.Is(err, errs.ErrSellerFlowNeedsWallet):
		return problems.SellerFlowNeedsWallet(errs.ErrSellerFlowNeedsWallet.Error())

	// providers and intents
	case errors.Is(err, errs.ErrProviderNotImplemented):
		return problems.ProviderNotImplemented(errs.ErrProviderNotImplemented.Error())
	case errors.Is(err, errs.ErrProviderRequestFailed):
		return problems.ProviderRequestFailed(errs.ErrProviderRequestFailed.Error())
	case errors.As(err, &unsupported):
		return problems.UnsupportedProvider(unsupported.Error(), unsupported.Provider)
	case errors.As(err, &provider):
		return providerProblem(provider)
	case errors.Is(err, errs.ErrSellerNotInWallet):
		return problems.SellerNotInWallet(errs.ErrSellerNotInWallet.Error())
	case errors.Is(err, errs.ErrNoTransactionToRefund):
		return problems.NoTransactionToRefund(errs.ErrNoTransactionToRefund.Error())
	case errors.Is(err, errs.ErrNoTransactionToCancel):
		return problems.NoTransactionToCancel(errs.ErrNoTransactionToCancel.Error())
	case errors.As(err, &intentStatus):
		return problems.IntentStatus(intentStatus.Error(), intentStatus.Action, intentStatus.Status)
	case errors.As(err, &revoked) && revoked.Op == "checkout":
		return problems.SellerRevokedCheckout(revoked.Error(), revoked.SellerID.String())
	case errors.As(err, &revoked):
		return problems.SellerRevoked(revoked.Error(), revoked.SellerID.String())
	case errors.As(err, &mismatch):
		return problems.SellerProviderMismatch(mismatch.Error(), mismatch.SellerID.String(), mismatch.Provider)

	// incoming webhooks
	case errors.Is(err, errs.ErrUnknownWebhookProvider):
		return problems.UnknownWebhookProvider(errs.ErrUnknownWebhookProvider.Error())
	case errors.Is(err, errs.ErrWebhookNotImplemented):
		return problems.WebhookNotImplemented(errs.ErrWebhookNotImplemented.Error())
	case errors.Is(err, errs.ErrWebhookPayloadInvalid):
		return problems.WebhookPayloadInvalid(errs.ErrWebhookPayloadInvalid.Error())
	case errors.Is(err, errs.ErrWebhookMissingDataID):
		return problems.WebhookMissingDataID(errs.ErrWebhookMissingDataID.Error())
	case errors.Is(err, errs.ErrWebhookSignatureMissing):
		return problems.WebhookSignatureMissing(errs.ErrWebhookSignatureMissing.Error())
	case errors.Is(err, errs.ErrWebhookSignatureMismatch):
		return problems.WebhookSignatureMismatch(errs.ErrWebhookSignatureMismatch.Error())
	case errors.As(err, &missingParam):
		return problems.Validation(missingParam.Error())
	case errors.As(err, &unhandled):
		return problems.UnhandledWebhookEvent(unhandled.Error(), unhandled.Type)
	case errors.As(err, &noIntent):
		return problems.WebhookIntentNotFound(noIntent.Error(), noIntent.PaymentID)
	case errors.As(err, &noToken):
		return problems.WebhookSellerTokenMissing(noToken.Error(), noToken.SellerID.String(), noToken.IntentID.String())

	// constraints
	case errors.Is(err, errs.ErrOrgWalletNameTaken):
		return problems.OrgWalletNameTaken(errs.ErrOrgWalletNameTaken.Error())
	case errors.Is(err, errs.ErrPersonalWalletTaken):
		return problems.PersonalWalletTaken(errs.ErrPersonalWalletTaken.Error())
	case errors.Is(err, errs.ErrDeliveryExists):
		return problems.DeliveryExists(errs.ErrDeliveryExists.Error())
	case errors.Is(err, errs.ErrWebhookEventExists):
		return problems.WebhookEventExists(errs.ErrWebhookEventExists.Error())
	case errors.Is(err, errs.ErrOrgSlugTaken):
		return problems.OrgSlugTaken(errs.ErrOrgSlugTaken.Error())
	case errors.Is(err, errs.ErrCredentialsConnected):
		return problems.CredentialsConnected(errs.ErrCredentialsConnected.Error())
	case errors.Is(err, errs.ErrOrgCollectorConnected):
		return problems.OrgCollectorConnected(errs.ErrOrgCollectorConnected.Error())
	case errors.Is(err, errs.ErrPersonalCollectorConnected):
		return problems.PersonalCollectorConnected(errs.ErrPersonalCollectorConnected.Error())
	case errors.Is(err, errs.ErrSellerConnected):
		return problems.SellerConnected(errs.ErrSellerConnected.Error())
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
