package commands

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"univents/internal/commerce/domain"
	"univents/internal/commerce/interfaces/http/dtos"
	"univents/internal/plataform/telemetry"
	"univents/internal/shared/sockets"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func mapPaymentError() string {
	return "payment could not be processed, please try again"
}

func (uc *CommandService) Purchase(ctx context.Context, conn *websocket.Conn, req dtos.BuyRequest, editionID uuid.UUID) error {
	user, err := uc.gaClient.Users.Get(ctx, req.UserID)
	if err != nil {
		return errors.New("unauthorized")
	}

	sessionID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(domain.ReservationDuration)

	// ── Phase 1: fetch + validate ─────────────────────────────────────────────
	ids := make([]uuid.UUID, len(req.Items))
	for i, item := range req.Items {
		ids[i] = item.ProductID
	}

	toBuy, err := uc.products.GetByIDs(ctx, ids)
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

	productMap := make(map[uuid.UUID]domain.Product, len(toBuy))
	for _, p := range toBuy {
		productMap[p.ID] = p
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
		_ = conn.WriteJSON(sockets.WSMessage{Type: "purchase_failed", Payload: invalid})
		return nil
	}

	for i, item := range req.Items {
		if p, ok := productMap[item.ProductID]; ok {
			req.Items[i].HasInventory = p.HasInventory
		}
	}

	// ── Phase 2: reserve ──────────────────────────────────────────────────────
	outcome, err := uc.products.ReserveItems(ctx, sessionID, req.Items, expiresAt)
	if err != nil {
		updates, uErr := uc.products.UnreserveItems(ctx, sessionID)
		if uErr != nil {
			telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
		}
		if len(updates) > 0 {
			_ = uc.inventory.Publish(ctx, editionID, updates)
		}
		return err
	}

	if len(outcome.InventoryUpdates) > 0 {
		_ = uc.inventory.Publish(ctx, editionID, outcome.InventoryUpdates)
	}

	for i, inv := range outcome.Unavailable {
		if p, ok := productMap[inv.ProductID]; ok {
			outcome.Unavailable[i].Name = p.Name
		}
	}

	if len(outcome.Reserved) == 0 {
		_ = conn.WriteJSON(sockets.WSMessage{
			Type:    "reservation_failed",
			Payload: map[string]any{"unavailable": outcome.Unavailable},
		})
		return nil
	}

	// build reservedDetails + total once, reused in partial_reservation and reservation_confirmed
	reservedDetails := make([]map[string]any, 0, len(outcome.Reserved))
	var total int
	for _, item := range outcome.Reserved {
		p := productMap[item.ProductID]
		total += p.PriceCents * item.Quantity
		reservedDetails = append(reservedDetails, map[string]any{
			"product_id":  p.ID,
			"name":        p.Name,
			"quantity":    item.Quantity,
			"price_cents": p.PriceCents,
		})
	}

	// ── Phase 3: partial confirmation ─────────────────────────────────────────
	if len(outcome.Unavailable) > 0 {
		confirmDeadline := time.Now().Add(60 * time.Second)

		_ = conn.WriteJSON(sockets.WSMessage{
			Type: "partial_reservation",
			Payload: map[string]any{
				"reserved":         reservedDetails,
				"unavailable":      outcome.Unavailable,
				"confirm_deadline": confirmDeadline,
			},
		})

		if err = conn.SetReadDeadline(confirmDeadline); err != nil {
			updates, uErr := uc.products.UnreserveItems(ctx, sessionID)
			if uErr != nil {
				telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
			}
			if len(updates) > 0 {
				_ = uc.inventory.Publish(ctx, editionID, updates)
			}
			return err
		}

		var confirmMsg sockets.WSMessage
		if err = conn.ReadJSON(&confirmMsg); err != nil {
			updates, uErr := uc.products.UnreserveItems(ctx, sessionID)
			if uErr != nil {
				telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
			}
			if len(updates) > 0 {
				_ = uc.inventory.Publish(ctx, editionID, updates)
			}
			return nil
		}

		if err = conn.SetReadDeadline(time.Time{}); err != nil {
			updates, uErr := uc.products.UnreserveItems(ctx, sessionID)
			if uErr != nil {
				telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
			}
			if len(updates) > 0 {
				_ = uc.inventory.Publish(ctx, editionID, updates)
			}
			return err
		}

		switch confirmMsg.Type {
		case "confirm_partial":
			// proceed with outcome.Reserved only
		case "cancel":
			updates, uErr := uc.products.UnreserveItems(ctx, sessionID)
			if uErr != nil {
				telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
			}
			if len(updates) > 0 {
				_ = uc.inventory.Publish(ctx, editionID, updates)
			}
			_ = conn.WriteJSON(sockets.WSMessage{Type: "reservation_cancelled"})
			return nil
		default:
			updates, uErr := uc.products.UnreserveItems(ctx, sessionID)
			if uErr != nil {
				telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
			}
			if len(updates) > 0 {
				_ = uc.inventory.Publish(ctx, editionID, updates)
			}
			_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "expected confirm_partial or cancel"})
			return nil
		}
	}

	// ── Phase 4: payment intent (no locks held) ───────────────────────────────
	intent, err := uc.payments.CreateIntent(ctx, paymentsSDK.CreateIntentRequest{
		Amount:   int64(total),
		Currency: "brl",
		Provider: viper.GetString("TRIEPAYMENTS_PROVIDER"),
		Metadata: json.RawMessage(`{"session_id": "` + sessionID.String() + `"}`),
	})
	if err != nil {
		updates, uErr := uc.products.UnreserveItems(ctx, sessionID)
		if uErr != nil {
			telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
		}
		if len(updates) > 0 {
			_ = uc.inventory.Publish(ctx, editionID, updates)
		}
		return err
	}
	paymentIntentID := intent.ID

	// ── Phase 5: purchase record TX (pure DB) ─────────────────────────────────
	if err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		pendingPurchase := domain.NewPurchase(domain.CreatePurchaseSpec{
			EditionID:       editionID,
			SessionID:       &sessionID,
			UserID:          user.ID,
			SubtotalCents:   int(intent.Amount),
			PaymentProvider: &intent.Provider,
			PaymentID:       &intent.ID,
		})

		purchase, err := uc.purchases.Create(ctx, *pendingPurchase)
		if err != nil {
			return err
		}

		for _, item := range outcome.Reserved {
			p := productMap[item.ProductID]

			if p.Type == domain.ProductTypeTicket {
				for range item.Quantity {
					if _, err = uc.purchases.CreateLineItem(ctx, domain.LineItem{
						PurchaseID:      purchase.ID,
						ItemType:        "ticket",
						ItemID:          *p.TicketID,
						Quantity:        1,
						UnitPriceCents:  p.PriceCents,
						TotalPriceCents: p.PriceCents,
					}); err != nil {
						return err
					}
				}
			} else {
				if _, err = uc.purchases.CreateLineItem(ctx, domain.LineItem{
					PurchaseID:      purchase.ID,
					ItemType:        "product",
					ItemID:          p.ID,
					Quantity:        item.Quantity,
					UnitPriceCents:  p.PriceCents,
					TotalPriceCents: p.PriceCents * item.Quantity,
				}); err != nil {
					return err
				}
			}
		}

		task, err := domain.NewReservationExpiredTask(sessionID, paymentIntentID, expiresAt, editionID)
		if err != nil {
			return err
		}
		if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
			_ = conn.WriteJSON(sockets.WSMessage{Type: "purchase_failed", Payload: "failed to initialize reservation timer"})
			return err
		}

		return nil
	}); err != nil {
		// FIXME: Call FailIntent when implemented
		updates, uErr := uc.products.UnreserveItems(ctx, sessionID)
		if uErr != nil {
			telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
		}
		if len(updates) > 0 {
			_ = uc.inventory.Publish(ctx, editionID, updates)
		}
		return err
	}

	// ── Phase 6: wait for payment ─────────────────────────────────────────────
	uc.ws.Register(sessionID.String(), conn)

	_ = conn.WriteJSON(sockets.WSMessage{
		Type: "reservation_confirmed",
		Payload: dtos.ReservationConfirmedPayload{
			SessionID: sessionID,
			ExpiresAt: expiresAt,
			Items:     reservedDetails,
			Total:     total,
		},
	})

	var payMsg sockets.WSMessage
	if err = conn.ReadJSON(&payMsg); err != nil {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "expected submit_payment"})
		return nil
	}

	if payMsg.Type != "submit_payment" {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "expected submit_payment message type"})
		return nil
	}

	payloadBytes, err := json.Marshal(payMsg.Payload)
	if err != nil {
		return err
	}

	var payReq dtos.SubmitPaymentPayload
	if err = json.Unmarshal(payloadBytes, &payReq); err != nil {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "invalid submit_payment payload"})
		return nil
	}

	_ = conn.WriteJSON(sockets.WSMessage{Type: "payment_processing"})

	if _, err = uc.payments.PayIntent(ctx, paymentIntentID, paymentsSDK.PayIntentRequest{
		CardToken:       payReq.CardToken,
		PaymentMethodID: payReq.PaymentMethodID,
		Installments:    payReq.Installments,
		PayerEmail:      user.Email,
	}); err != nil {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "payment_failed", Payload: map[string]string{"reason": mapPaymentError()}})
		return nil
	}

	// block with timeout waiting for webhook to resolve via ws.Notify
	paymentTimeout := time.After(30 * time.Second)
	connClosed := make(chan struct{})

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				close(connClosed)
				return
			}
		}
	}()

	select {
	case <-paymentTimeout:
		_ = conn.WriteJSON(sockets.WSMessage{
			Type:    "payment_pending",
			Payload: "payment is taking longer than expected, you can close this and check your purchases",
		})
		return nil
	case <-connClosed:
		return nil
	}
}
