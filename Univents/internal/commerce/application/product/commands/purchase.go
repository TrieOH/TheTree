package commands

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"univents/internal/commerce/domain"
	"univents/internal/commerce/interfaces/http/dtos"
	"univents/internal/shared/authz"
	"univents/internal/shared/sockets"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
)

var ErrReservationFailed = errors.New("reservation_failed")

func (uc *CommandService) Purchase(ctx context.Context, conn *websocket.Conn, req dtos.BuyRequest, editionID uuid.UUID) (err error) {
	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(domain.ReservationDuration)
	sessionID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	var paymentIntentID string

	if err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		ids := make([]uuid.UUID, len(req.Items))
		quantityMap := make(map[uuid.UUID]int, len(req.Items))
		for i, item := range req.Items {
			ids[i] = item.ProductID
			quantityMap[item.ProductID] = item.Quantity
		}

		var toBuy []domain.Product
		toBuy, err = uc.products.GetByIDs(ctx, ids)
		if err != nil {
			return err
		}

		if len(toBuy) != len(ids) {
			_ = conn.WriteJSON(sockets.WSMessage{
				Type: "purchase_failed",
				Payload: map[string]any{
					"reason":      "one or more products do not exist",
					"product_ids": ids,
				},
			})
			return nil
		}

		invalid := make([]domain.InvalidProduct, 0)

		for _, p := range toBuy {
			if p.Status != domain.ProductStatusAvailable {
				var reason string

				switch p.Status {
				case domain.ProductStatusDraft:
					reason = "product is not yet available"
				case domain.ProductStatusSoldOut:
					reason = "product is sold out"
				case domain.ProductStatusUnavailable:
					reason = "product is unavailable"
				default:
					reason = "product cannot be purchased"
				}

				invalid = append(invalid, domain.InvalidProduct{
					ProductID: p.ID,
					Name:      p.Name,
					Reason:    reason,
				})
			}
		}

		if len(invalid) > 0 {
			_ = conn.WriteJSON(sockets.WSMessage{
				Type:    "purchase_failed",
				Payload: invalid,
			})
			return nil
		}

		// enrich cart items with inventory flag
		productMap := make(map[uuid.UUID]domain.Product, len(toBuy))
		for _, p := range toBuy {
			productMap[p.ID] = p
		}
		for i, item := range req.Items {
			if p, ok := productMap[item.ProductID]; ok {
				req.Items[i].HasInventory = p.HasInventory
			}
		}

		if err := uc.products.ReserveItems(ctx, sessionID, req.Items, expiresAt); err != nil {
			return ErrReservationFailed
		}

		var subtotal int
		for _, p := range toBuy {
			subtotal += p.PriceCents * quantityMap[p.ID]
		}

		var intent *paymentsSDK.Intent
		intent, err = uc.payments.CreateIntent(ctx, paymentsSDK.CreateIntentRequest{
			Amount:   int64(subtotal),
			Currency: "brl",
			Provider: "mock",
			Metadata: json.RawMessage(`{"session_id": "` + sessionID.String() + `"}`),
		})
		if err != nil {
			return err
		}

		paymentIntentID = intent.ID

		pendingPurchase := domain.NewPurchase(domain.CreatePurchaseSpec{
			EditionID:       editionID,
			SessionID:       &sessionID,
			UserID:          sub.ID,
			SubtotalCents:   int(intent.Amount),
			PaymentProvider: &intent.Provider,
			PaymentID:       &intent.ID,
		})

		var purchase *domain.Purchase
		purchase, err = uc.purchases.Create(ctx, *pendingPurchase)
		if err != nil {
			return err
		}

		for _, product := range toBuy {
			quantity := quantityMap[product.ID]

			if product.Type == domain.ProductTypeTicket {
				for range quantity {
					_, err = uc.purchases.CreateLineItem(ctx, domain.LineItem{
						PurchaseID:      purchase.ID,
						ItemType:        "ticket",
						ItemID:          *product.TicketID,
						Quantity:        1,
						UnitPriceCents:  product.PriceCents,
						TotalPriceCents: product.PriceCents,
					})
					if err != nil {
						return err
					}
				}
			} else {
				_, err = uc.purchases.CreateLineItem(ctx, domain.LineItem{
					PurchaseID:      purchase.ID,
					ItemType:        "product",
					ItemID:          product.ID,
					Quantity:        quantity,
					UnitPriceCents:  product.PriceCents,
					TotalPriceCents: product.PriceCents * quantity,
				})
				if err != nil {
					return err
				}
			}
		}

		var task *asynq.Task
		task, err = domain.NewReservationExpiredTask(sessionID, paymentIntentID, expiresAt)
		if err != nil {
			return err
		}
		if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
			_ = conn.WriteJSON(sockets.WSMessage{
				Type:    "purchase_failed",
				Payload: "failed to initialize reservation timer, please try again",
			})
			return err
		}

		return nil
	}); err != nil {
		if errors.Is(err, ErrReservationFailed) {
			_ = conn.WriteJSON(sockets.WSMessage{
				Type:    "reservation_failed",
				Payload: "one or more items are unavailable or already reserved",
			})
			return nil
		}
		return err
	}

	uc.ws.Register(sessionID.String(), conn)

	_ = conn.WriteJSON(sockets.WSMessage{
		Type: "reservation_confirmed",
		Payload: dtos.ReservationConfirmedPayload{
			SessionID: sessionID,
			ExpiresAt: expiresAt,
		},
	})

	return nil
}
