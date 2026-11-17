package errs

import (
	"fmt"
)

// PaymentRejectedError is a payment Payssage refused for something the
// buyer can fix (a rejected card, an invalid token).
type PaymentRejectedError struct{ Message string }

func (e *PaymentRejectedError) Error() string { return "payment failed: " + e.Message }

// PurchaseStatusError is a purchase whose status does not allow a refund.
type PurchaseStatusError struct{ Status string }

func (e *PurchaseStatusError) Error() string {
	return fmt.Sprintf("purchase cannot be refunded from status %q", e.Status)
}
