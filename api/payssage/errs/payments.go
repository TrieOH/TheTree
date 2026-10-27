package errs

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Providers.
var (
	ErrProviderNotImplemented = errors.New("provider not implemented")
	// ErrProviderRequestFailed is a provider call that never got an answer
	// (network, timeout); the transport error is wrapped for the logs.
	ErrProviderRequestFailed = errors.New("payment provider request failed")
)

// UnsupportedProviderError is a provider name Payssage does not know.
type UnsupportedProviderError struct{ Provider string }

func (e *UnsupportedProviderError) Error() string { return "unsupported provider " + e.Provider }

// ProviderError is a provider's own error answer. Status is the HTTP status
// the provider answered with; the mapper derives Payssage's status from it.
type ProviderError struct {
	Provider         string
	Status           int
	Message          string
	Code             string
	CauseCode        int
	CauseDescription string
	// Reason is Payssage's actionable reading of a known cause code.
	Reason string
}

func (e *ProviderError) Error() string {
	if e.Reason != "" {
		return e.Provider + ": " + e.Message + " (" + e.Reason + ")"
	}
	return e.Provider + ": " + e.Message
}

// Intents.
var (
	ErrSellerNotInWallet     = errors.New("seller does not belong to wallet")
	ErrNoTransactionToRefund = errors.New("intent has no mercadopago transaction to refund")
	ErrNoTransactionToCancel = errors.New("intent has no mercadopago transaction to cancel")
)

// IntentStatusError is an intent whose status does not allow the action
// ("refunded", "cancelled").
type IntentStatusError struct {
	Action string
	Status string
}

func (e *IntentStatusError) Error() string {
	return fmt.Sprintf("intent cannot be %s from status %q", e.Action, e.Status)
}

// SellerRevokedError is a provider operation on a seller that revoked
// Payssage's access. Op is the operation: "checkout", "refund", "cancel".
type SellerRevokedError struct {
	SellerID uuid.UUID
	Op       string
}

func (e *SellerRevokedError) Error() string {
	switch e.Op {
	case "refund":
		return fmt.Sprintf("seller %s has revoked mercadopago access; cannot refund — payment may require manual handling at mercadopago", e.SellerID)
	case "cancel":
		return fmt.Sprintf("seller %s has revoked mercadopago access; cannot cancel — payment may still be live at mercadopago and requires manual review", e.SellerID)
	}
	return fmt.Sprintf("seller %s has revoked provider access", e.SellerID)
}

// SellerProviderMismatchError is a seller connected through another provider.
type SellerProviderMismatchError struct {
	SellerID uuid.UUID
	Provider string
}

func (e *SellerProviderMismatchError) Error() string {
	return fmt.Sprintf("seller %s is not a mercadopago seller (got %q)", e.SellerID, e.Provider)
}
