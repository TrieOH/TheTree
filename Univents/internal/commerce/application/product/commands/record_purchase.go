package commands

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/errx"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
)

type recordPurchaseInput struct {
	session *domain.PurchaseSession
	intent  *paymentsSDK.Intent
}

func (uc *CommandService) recordPurchase(ctx context.Context, in recordPurchaseInput) error {
	if err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		pendingPurchase := domain.NewPurchase(domain.CreatePurchaseSpec{
			EditionID:       in.session.EditionID,
			SessionID:       &in.session.SessionID,
			UserID:          in.session.UserID,
			SubtotalCents:   int(in.intent.Amount),
			PaymentProvider: &in.intent.Provider,
			PaymentID:       &in.intent.ID,
		})

		purchase, err := uc.purchases.Create(ctx, *pendingPurchase)
		if err != nil {
			if errx.IsKind(err, "not_found") {
				// ON CONFLICT DO NOTHING on session_id returns no rows — purchase already recorded, safe to skip
				return nil
			}
			return err
		}

		for _, item := range in.session.Reserved {
			if item.ProductType == domain.ProductTypeTicket {
				for range item.Quantity {
					if _, err = uc.purchases.CreateLineItem(ctx, domain.LineItem{
						PurchaseID:      purchase.ID,
						ItemType:        "ticket",
						ItemID:          *item.TicketID,
						Quantity:        1,
						UnitPriceCents:  item.PriceCents,
						TotalPriceCents: item.PriceCents,
					}); err != nil {
						return err
					}
				}
			} else {
				if _, err = uc.purchases.CreateLineItem(ctx, domain.LineItem{
					PurchaseID:      purchase.ID,
					ItemType:        "product",
					ItemID:          item.ProductID,
					Quantity:        item.Quantity,
					UnitPriceCents:  item.PriceCents,
					TotalPriceCents: item.PriceCents * item.Quantity,
				}); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
