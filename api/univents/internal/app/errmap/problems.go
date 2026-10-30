package errmap

import (
	"errors"

	"univents/errs"
	"univents/internal/problems"

	fp "github.com/MintzyG/problematic"
)

// Problem answers Univents's named errors with their declared problems
// (ADR-0019), and nil for an error it doesn't know; problems.Mapper then
// tries the shared cases. Named errors don't wrap named errors, so case
// order never breaks a tie. The detail is the named error's message, never
// the chain's.
func Problem(err error) *fp.Problem {
	var (
		rejected       *errs.PaymentRejectedError
		purchaseStatus *errs.PurchaseStatusError
		accessLevel    *errs.AccessLevelError
	)
	switch {
	// request shapes
	case errors.Is(err, errs.ErrInvalidRequest):
		return problems.InvalidRequest(errs.ErrInvalidRequest.Error())
	case errors.Is(err, errs.ErrInvalidDesignData):
		return problems.InvalidDesignData(errs.ErrInvalidDesignData.Error())
	case errors.Is(err, errs.ErrInvalidGalleryURLs):
		return problems.InvalidGalleryURLs(errs.ErrInvalidGalleryURLs.Error())

	// events, editions, templates, payments
	case errors.Is(err, errs.ErrEventNotDraft):
		return problems.EventNotDraft(errs.ErrEventNotDraft.Error())
	case errors.Is(err, errs.ErrEventNotActive):
		return problems.EventNotActive(errs.ErrEventNotActive.Error())
	case errors.Is(err, errs.ErrEditionAlreadyPublished):
		return problems.EditionAlreadyPublished(errs.ErrEditionAlreadyPublished.Error())
	case errors.Is(err, errs.ErrInvalidBadgeOrigin):
		return problems.InvalidBadgeOrigin(errs.ErrInvalidBadgeOrigin.Error())
	case errors.Is(err, errs.ErrStaffTemplateWithTicketType):
		return problems.StaffTemplateWithTicketType(errs.ErrStaffTemplateWithTicketType.Error())
	case errors.Is(err, errs.ErrSellerNotInPlatformWallet):
		return problems.SellerNotInPlatformWallet(errs.ErrSellerNotInPlatformWallet.Error())
	case errors.Is(err, errs.ErrInvalidPayssageWebhookSignature):
		return problems.InvalidPayssageWebhookSignature(errs.ErrInvalidPayssageWebhookSignature.Error())

	// checkout
	case errors.Is(err, errs.ErrEditionNotPublished):
		return problems.EditionNotPublished(errs.ErrEditionNotPublished.Error())
	case errors.Is(err, errs.ErrNoPaymentConfig):
		return problems.NoPaymentConfig(errs.ErrNoPaymentConfig.Error())
	case errors.Is(err, errs.ErrNoItems):
		return problems.NoItems(errs.ErrNoItems.Error())
	case errors.Is(err, errs.ErrQuantityBelowOne):
		return problems.QuantityBelowOne(errs.ErrQuantityBelowOne.Error())
	case errors.Is(err, errs.ErrTicketQuantity):
		return problems.TicketQuantity(errs.ErrTicketQuantity.Error())
	case errors.Is(err, errs.ErrTicketNeedsAttendee):
		return problems.TicketNeedsAttendee(errs.ErrTicketNeedsAttendee.Error())
	case errors.Is(err, errs.ErrTicketNeedsAttendeeEmail):
		return problems.TicketNeedsAttendeeEmail(errs.ErrTicketNeedsAttendeeEmail.Error())
	case errors.Is(err, errs.ErrOneTicketPerPerson):
		return problems.OneTicketPerPerson(errs.ErrOneTicketPerPerson.Error())
	case errors.Is(err, errs.ErrDuplicateProductLine):
		return problems.DuplicateProductLine(errs.ErrDuplicateProductLine.Error())
	case errors.Is(err, errs.ErrProgramQuantity):
		return problems.ProgramQuantity(errs.ErrProgramQuantity.Error())
	case errors.Is(err, errs.ErrDuplicateProgramLine):
		return problems.DuplicateProgramLine(errs.ErrDuplicateProgramLine.Error())
	case errors.Is(err, errs.ErrInvalidItemType):
		return problems.InvalidItemType(errs.ErrInvalidItemType.Error())
	case errors.Is(err, errs.ErrProgramNeedsTicket):
		return problems.ProgramNeedsTicket(errs.ErrProgramNeedsTicket.Error())
	case errors.Is(err, errs.ErrProgramNeedsTicketInCart):
		return problems.ProgramNeedsTicketInCart(errs.ErrProgramNeedsTicketInCart.Error())
	case errors.Is(err, errs.ErrAttendeeEmailUnknown):
		return problems.AttendeeEmailUnknown(errs.ErrAttendeeEmailUnknown.Error())
	case errors.Is(err, errs.ErrAttendeeUserMismatch):
		return problems.AttendeeUserMismatch(errs.ErrAttendeeUserMismatch.Error())
	case errors.Is(err, errs.ErrCardTokenRequired):
		return problems.CardTokenRequired(errs.ErrCardTokenRequired.Error())
	case errors.Is(err, errs.ErrPaymentMethodIDRequired):
		return problems.PaymentMethodIDRequired(errs.ErrPaymentMethodIDRequired.Error())
	case errors.Is(err, errs.ErrInstallmentsBelowOne):
		return problems.InstallmentsBelowOne(errs.ErrInstallmentsBelowOne.Error())
	case errors.Is(err, errs.ErrInvalidPaymentMethod):
		return problems.InvalidPaymentMethod(errs.ErrInvalidPaymentMethod.Error())
	case errors.Is(err, errs.ErrUnknownItem):
		return problems.UnknownItem(errs.ErrUnknownItem.Error())
	case errors.Is(err, errs.ErrItemNotInEdition):
		return problems.ItemNotInEdition(errs.ErrItemNotInEdition.Error())
	case errors.Is(err, errs.ErrPersonHasTicket):
		return problems.PersonHasTicket(errs.ErrPersonHasTicket.Error())
	case errors.Is(err, errs.ErrItemsUnavailable):
		return problems.ItemsUnavailable(errs.ErrItemsUnavailable.Error())
	case errors.As(err, &rejected):
		return problems.PaymentRejected(rejected.Error(), rejected.Message)

	// refunds
	case errors.Is(err, errs.ErrAlreadyRefunded):
		return problems.AlreadyRefunded(errs.ErrAlreadyRefunded.Error())
	case errors.Is(err, errs.ErrFreeOrderRefund):
		return problems.FreeOrderRefund(errs.ErrFreeOrderRefund.Error())
	case errors.As(err, &purchaseStatus):
		return problems.PurchaseStatus(purchaseStatus.Error(), purchaseStatus.Status)

	// programs
	case errors.Is(err, errs.ErrTicketRequiredForActivity):
		return problems.TicketRequiredForActivity(errs.ErrTicketRequiredForActivity.Error())
	case errors.Is(err, errs.ErrTicketRequiredForCheckpoint):
		return problems.TicketRequiredForCheckpoint(errs.ErrTicketRequiredForCheckpoint.Error())
	case errors.As(err, &accessLevel):
		return problems.AccessLevel(accessLevel.Error(), accessLevel.Kind)
	case errors.Is(err, errs.ErrAlreadyRegistered):
		return problems.AlreadyRegistered(errs.ErrAlreadyRegistered.Error())
	case errors.Is(err, errs.ErrOccurrenceFull):
		return problems.OccurrenceFull(errs.ErrOccurrenceFull.Error())
	case errors.Is(err, errs.ErrNotRegistered):
		return problems.NotRegistered(errs.ErrNotRegistered.Error())
	case errors.Is(err, errs.ErrRegistrationLocked):
		return problems.RegistrationLocked(errs.ErrRegistrationLocked.Error())
	case errors.Is(err, errs.ErrCheckInNotCheckpoint):
		return problems.CheckInNotCheckpoint(errs.ErrCheckInNotCheckpoint.Error())
	case errors.Is(err, errs.ErrParticipationCancelled):
		return problems.ParticipationCancelled(errs.ErrParticipationCancelled.Error())

	// signatures
	case errors.Is(err, errs.ErrInvalidSignatureToken):
		return problems.InvalidSignatureToken(errs.ErrInvalidSignatureToken.Error())
	case errors.Is(err, errs.ErrSignatureRequestNotPending):
		return problems.SignatureRequestNotPending(errs.ErrSignatureRequestNotPending.Error())
	case errors.Is(err, errs.ErrTokenRequestEditionMismatch):
		return problems.TokenRequestEditionMismatch(errs.ErrTokenRequestEditionMismatch.Error())
	case errors.Is(err, errs.ErrTokenSignatureEditionMismatch):
		return problems.TokenSignatureEditionMismatch(errs.ErrTokenSignatureEditionMismatch.Error())
	case errors.Is(err, errs.ErrSignatureRequestHandled):
		return problems.SignatureRequestHandled(errs.ErrSignatureRequestHandled.Error())
	case errors.Is(err, errs.ErrSignatureRequestCompleted):
		return problems.SignatureRequestCompleted(errs.ErrSignatureRequestCompleted.Error())

	// constraints
	case errors.Is(err, errs.ErrEditionOverlap):
		return problems.EditionOverlap(errs.ErrEditionOverlap.Error())
	case errors.Is(err, errs.ErrTicketHeld):
		return problems.TicketHeld(errs.ErrTicketHeld.Error())
	case errors.Is(err, errs.ErrProductCodeTaken):
		return problems.ProductCodeTaken(errs.ErrProductCodeTaken.Error())
	case errors.Is(err, errs.ErrVariantCodeTaken):
		return problems.VariantCodeTaken(errs.ErrVariantCodeTaken.Error())
	case errors.Is(err, errs.ErrOccurrenceRegistrationTaken):
		return problems.OccurrenceRegistrationTaken(errs.ErrOccurrenceRegistrationTaken.Error())
	case errors.Is(err, errs.ErrIdempotencyKeyUsed):
		return problems.IdempotencyKeyUsed(errs.ErrIdempotencyKeyUsed.Error())
	case errors.Is(err, errs.ErrBadgeEmitted):
		return problems.BadgeEmitted(errs.ErrBadgeEmitted.Error())
	case errors.Is(err, errs.ErrBadgeTemplateScopeTaken):
		return problems.BadgeTemplateScopeTaken(errs.ErrBadgeTemplateScopeTaken.Error())
	case errors.Is(err, errs.ErrPendingPurchaseExists):
		return problems.PendingPurchaseExists(errs.ErrPendingPurchaseExists.Error())
	case errors.Is(err, errs.ErrTicketUnitInPurchase):
		return problems.TicketUnitInPurchase(errs.ErrTicketUnitInPurchase.Error())
	case errors.Is(err, errs.ErrProductLineInPurchase):
		return problems.ProductLineInPurchase(errs.ErrProductLineInPurchase.Error())
	case errors.Is(err, errs.ErrProgramLineInPurchase):
		return problems.ProgramLineInPurchase(errs.ErrProgramLineInPurchase.Error())
	case errors.Is(err, errs.ErrWSTokenIssued):
		return problems.WSTokenIssued(errs.ErrWSTokenIssued.Error())
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
	errs.ErrEditionNotFound,
	errs.ErrBadgeTemplateNotFound,
	errs.ErrProgramNotFound,
	errs.ErrPurchaseNotFound,
	errs.ErrTicketTypeNotFound,
	errs.ErrRegistrationNotFound,
	errs.ErrCertificationNotFound,
	errs.ErrEventNotFound,
	errs.ErrProductNotFound,
	errs.ErrSignatureNotFound,
	errs.ErrWSTokenNotFound,
}

// invalid are the check constraints, answered as the shared validation
// problem.
var invalid = []error{
	errs.ErrInvalidEventStatus,
	errs.ErrPublicKeyNeedsSeller,
	errs.ErrInvalidEventMemberRole,
	errs.ErrEditionDates,
	errs.ErrRegistrationAfterStart,
	errs.ErrInvalidRegistrationStatus,
	errs.ErrInvalidProductPurchaseStatus,
	errs.ErrInvalidProgramKind,
	errs.ErrOccurrenceDates,
	errs.ErrInvalidParticipationStatus,
	errs.ErrInvalidSignatureRequestStatus,
	errs.ErrInvalidSignatureStatus,
	errs.ErrReadySignatureNeedsImage,
	errs.ErrInvalidCertificationKind,
	errs.ErrInvalidBadgeEmissionOrigin,
	errs.ErrInvalidBadgeEmissionStatus,
	errs.ErrInvalidBadgeTemplateScope,
	errs.ErrInvalidPurchaseStatus,
	errs.ErrInvalidPurchaseItemType,
	errs.ErrPurchaseItemQuantity,
}
