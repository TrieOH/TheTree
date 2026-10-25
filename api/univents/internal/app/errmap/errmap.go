// Package errmap is Univents' one table from named errors to HTTP
// (ADR-0016). It sits apart from internal/app, which wires the features,
// so service tests can check what status an error answers with without
// importing the features they belong to.
package errmap

import (
	"net/http"

	"lib/httpserver"
	"univents/errs"
)

// Errors is the table. The most specific error in a chain wins, and an
// error missing here is a generic 500. Request fields attached with
// lib/fielderr are listed in the response. The shared lib errors (missing
// rows, constraint classes, batches, insufficient role) are mapped by
// lib/httpserver after these.
func Errors() httpserver.ErrorMap {
	const (
		badRequest = http.StatusBadRequest
		forbidden  = http.StatusForbidden
		notFound   = http.StatusNotFound
		conflict   = http.StatusConflict
	)
	is := httpserver.Is
	invalid := func(target error) httpserver.ErrorMapping { return is(target, badRequest).Validation() }

	return httpserver.ErrorMap{
		// request shapes
		is(errs.ErrInvalidRequest, badRequest),
		is(errs.ErrInvalidDesignData, badRequest),
		is(errs.ErrInvalidGalleryURLs, badRequest),

		// events, editions, templates, payments
		is(errs.ErrEventNotDraft, badRequest),
		is(errs.ErrEventNotActive, badRequest),
		is(errs.ErrEditionAlreadyPublished, badRequest),
		is(errs.ErrInvalidBadgeOrigin, badRequest),
		is(errs.ErrStaffTemplateWithTicketType, badRequest),
		is(errs.ErrSellerNotInPlatformWallet, badRequest),
		is(errs.ErrInvalidPayssageWebhookSignature, badRequest),

		// checkout
		is(errs.ErrEditionNotPublished, badRequest),
		is(errs.ErrNoPaymentConfig, badRequest),
		is(errs.ErrNoItems, badRequest),
		is(errs.ErrQuantityBelowOne, badRequest),
		is(errs.ErrTicketQuantity, badRequest),
		is(errs.ErrTicketNeedsAttendee, badRequest),
		is(errs.ErrTicketNeedsAttendeeEmail, badRequest),
		is(errs.ErrOneTicketPerPerson, badRequest),
		is(errs.ErrDuplicateProductLine, badRequest),
		is(errs.ErrProgramQuantity, badRequest),
		is(errs.ErrDuplicateProgramLine, badRequest),
		is(errs.ErrInvalidItemType, badRequest),
		is(errs.ErrProgramNeedsTicket, badRequest),
		is(errs.ErrProgramNeedsTicketInCart, badRequest),
		is(errs.ErrAttendeeEmailUnknown, badRequest),
		is(errs.ErrAttendeeUserMismatch, badRequest),
		is(errs.ErrCardTokenRequired, badRequest),
		is(errs.ErrPaymentMethodIDRequired, badRequest),
		is(errs.ErrInstallmentsBelowOne, badRequest),
		is(errs.ErrInvalidPaymentMethod, badRequest),
		is(errs.ErrUnknownItem, badRequest),
		is(errs.ErrItemNotInEdition, badRequest),
		is(errs.ErrPersonHasTicket, conflict),
		is(errs.ErrItemsUnavailable, conflict),
		httpserver.As[*errs.PaymentRejectedError](badRequest),

		// refunds
		is(errs.ErrAlreadyRefunded, conflict),
		is(errs.ErrFreeOrderRefund, badRequest),
		httpserver.As[*errs.PurchaseStatusError](conflict),

		// programs
		is(errs.ErrTicketRequiredForActivity, forbidden),
		is(errs.ErrTicketRequiredForCheckpoint, forbidden),
		httpserver.As[*errs.AccessLevelError](forbidden),
		is(errs.ErrAlreadyRegistered, conflict),
		is(errs.ErrOccurrenceFull, conflict),
		is(errs.ErrNotRegistered, notFound),
		is(errs.ErrRegistrationLocked, conflict),
		is(errs.ErrCheckInNotCheckpoint, badRequest),
		is(errs.ErrParticipationCancelled, conflict),

		// signatures
		is(errs.ErrInvalidSignatureToken, forbidden),
		is(errs.ErrSignatureRequestNotPending, conflict),
		is(errs.ErrTokenRequestEditionMismatch, forbidden),
		is(errs.ErrTokenSignatureEditionMismatch, forbidden),
		is(errs.ErrSignatureRequestHandled, notFound),
		is(errs.ErrSignatureRequestCompleted, notFound),

		// not found
		is(errs.ErrEditionNotFound, notFound),
		is(errs.ErrBadgeTemplateNotFound, notFound),
		is(errs.ErrProgramNotFound, notFound),
		is(errs.ErrPurchaseNotFound, notFound),
		is(errs.ErrTicketTypeNotFound, notFound),
		is(errs.ErrRegistrationNotFound, notFound),
		is(errs.ErrCertificationNotFound, notFound),
		is(errs.ErrEventNotFound, notFound),
		is(errs.ErrProductNotFound, notFound),
		is(errs.ErrSignatureNotFound, notFound),
		is(errs.ErrWSTokenNotFound, notFound),

		// constraints: unique and exclusion → 409, check → 400 validation
		invalid(errs.ErrInvalidEventStatus),
		invalid(errs.ErrPublicKeyNeedsSeller),
		invalid(errs.ErrInvalidEventMemberRole),
		invalid(errs.ErrEditionDates),
		invalid(errs.ErrRegistrationAfterStart),
		is(errs.ErrEditionOverlap, conflict),
		invalid(errs.ErrInvalidRegistrationStatus),
		is(errs.ErrTicketHeld, conflict),
		is(errs.ErrProductCodeTaken, conflict),
		is(errs.ErrVariantCodeTaken, conflict),
		invalid(errs.ErrInvalidProductPurchaseStatus),
		invalid(errs.ErrInvalidProgramKind),
		invalid(errs.ErrOccurrenceDates),
		invalid(errs.ErrInvalidParticipationStatus),
		is(errs.ErrOccurrenceRegistrationTaken, conflict),
		invalid(errs.ErrInvalidSignatureRequestStatus),
		is(errs.ErrIdempotencyKeyUsed, conflict),
		invalid(errs.ErrInvalidSignatureStatus),
		invalid(errs.ErrReadySignatureNeedsImage),
		invalid(errs.ErrInvalidCertificationKind),
		invalid(errs.ErrInvalidBadgeEmissionOrigin),
		invalid(errs.ErrInvalidBadgeEmissionStatus),
		is(errs.ErrBadgeEmitted, conflict),
		invalid(errs.ErrInvalidBadgeTemplateScope),
		is(errs.ErrBadgeTemplateScopeTaken, conflict),
		invalid(errs.ErrInvalidPurchaseStatus),
		is(errs.ErrPendingPurchaseExists, conflict),
		invalid(errs.ErrInvalidPurchaseItemType),
		invalid(errs.ErrPurchaseItemQuantity),
		is(errs.ErrTicketUnitInPurchase, conflict),
		is(errs.ErrProductLineInPurchase, conflict),
		is(errs.ErrProgramLineInPurchase, conflict),
		is(errs.ErrWSTokenIssued, conflict),
	}
}
