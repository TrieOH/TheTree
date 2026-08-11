// Package purchases serves the buyer's purchase history read (split 5):
// the authenticated user's purchases with their items, newest first.
// Strictly read-only — the webhook receiver (split 4) and the expiry
// worker (split 7) own status changes.
package purchases

import (
	"context"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
)

type Operations struct {
	purchases ports.PurchaseRepo
}

func NewOperations(purchases ports.PurchaseRepo) *Operations {
	return &Operations{purchases: purchases}
}

// PurchaseDetail is one purchase with its items — the list entry shape.
type PurchaseDetail struct {
	Purchase models.Purchase
	Items    []models.PurchaseItem
}

// ListForUser returns the user's purchases newest first (the repo orders by
// created_at DESC), each with its items from the availability ledger (D4).
// No pagination yet (decision taken at review): the list is unbounded for
// now and pages later.
func (o *Operations) ListForUser(ctx context.Context, userID uuid.UUID) ([]PurchaseDetail, error) {
	list, err := o.purchases.ListByPurchaser(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]PurchaseDetail, 0, len(list))
	for _, purchase := range list {
		items, err := o.purchases.ListItemsByPurchase(ctx, purchase.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, PurchaseDetail{Purchase: purchase, Items: items})
	}
	return out, nil
}
