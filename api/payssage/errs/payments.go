package errs

import (
	"fmt"
	"uuid"
)

// ProviderError is a provider's own error answer. Status is the HTTP status
// the provider answered with; the mapper derives Payssage's status from it,
// so it stays by hand: one error answers several problems.
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

// SellerRevokedError is a provider operation on a seller that revoked
// Payssage's access. Op is the operation: "checkout", "refund", "cancel".
// It stays by hand: Op picks the problem and the message.
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
