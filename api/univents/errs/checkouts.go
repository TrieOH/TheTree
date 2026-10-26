package errs

import (
	"errors"
	"fmt"
)

// Checkout preconditions.
var (
	ErrEditionNotPublished = errors.New("edition is not published")
	ErrNoPaymentConfig     = errors.New("event has no payment config")
)

// Checkout cart and payment validation; returned with the offending fields
// through lib/fielderr.
var (
	ErrNoItems                  = errors.New("checkout requires at least one item")
	ErrQuantityBelowOne         = errors.New("item quantity must be at least 1")
	ErrTicketQuantity           = errors.New("ticket lines are one per person (quantity must be 1)")
	ErrTicketNeedsAttendee      = errors.New("ticket lines require an attendee")
	ErrTicketNeedsAttendeeEmail = errors.New("ticket lines require an attendee email")
	ErrOneTicketPerPerson       = errors.New("only one ticket per person")
	ErrDuplicateProductLine     = errors.New("duplicate product line")
	ErrProgramQuantity          = errors.New("program lines are one session per purchase (quantity must be 1)")
	ErrDuplicateProgramLine     = errors.New("duplicate program line")
	ErrInvalidItemType          = errors.New("invalid item_type")
	ErrProgramNeedsTicket       = errors.New("program items require a ticket item")
	ErrProgramNeedsTicketInCart = errors.New("program items require a ticket item in the same cart")
	ErrAttendeeEmailUnknown     = errors.New("attendee email does not match an account")
	ErrAttendeeUserMismatch     = errors.New("attendee user_id does not match the email")
	ErrCardTokenRequired        = errors.New("card_token is required for card payments")
	ErrPaymentMethodIDRequired  = errors.New("payment_method_id is required for card payments")
	ErrInstallmentsBelowOne     = errors.New("installments must be at least 1")
	ErrInvalidPaymentMethod     = errors.New("invalid payment_method")
	ErrUnknownItem              = errors.New("unknown item")
	ErrItemNotInEdition         = errors.New("item does not belong to this edition")
	ErrPersonHasTicket          = errors.New("this person already has a ticket for this edition")
	ErrItemsUnavailable         = errors.New("some items are no longer available")
)

// PaymentRejectedError is a payment Payssage refused for something the
// buyer can fix (a rejected card, an invalid token).
type PaymentRejectedError struct{ Message string }

func (e *PaymentRejectedError) Error() string { return "payment failed: " + e.Message }

// Refunds.
var (
	ErrAlreadyRefunded = errors.New("purchase already refunded")
	ErrFreeOrderRefund = errors.New("purchase has no payment intent (free order) and cannot be refunded")
)

// PurchaseStatusError is a purchase whose status does not allow a refund.
type PurchaseStatusError struct{ Status string }

func (e *PurchaseStatusError) Error() string {
	return fmt.Sprintf("purchase cannot be refunded from status %q", e.Status)
}
