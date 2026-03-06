package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"univents/internal/commerce/domain"
	"univents/internal/commerce/interfaces/http/dtos"
	"univents/internal/shared/authz"
	"univents/internal/shared/sockets"

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
	var paymentIntentID, clientSecret, provider string

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

		paymentIntentID, clientSecret, provider, err = uc.payments.CreatePaymentIntent(ctx, req)
		if err != nil {
			return err
		}

		pendingPurchase := domain.NewPurchase(domain.CreatePurchaseSpec{
			EditionID:       editionID,
			UserID:          sub.ID,
			SubtotalCents:   subtotal,
			PaymentProvider: &provider,
			PaymentID:       &paymentIntentID,
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
			SessionID:    sessionID,
			ClientSecret: clientSecret,
			ExpiresAt:    expiresAt,
		},
	})

	return nil
}

func (uc *CommandService) ConfirmPayment(ctx context.Context, sessionID uuid.UUID, paymentIntentID string) error {
	// 1. mark items as sold in db
	if err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := uc.products.DeleteReservation(ctx, sessionID); err != nil {
			return err
		}
		if err := uc.purchases.ConfirmPurchase(ctx, paymentIntentID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err := uc.ws.Notify(sessionID.String(), sockets.WSMessage{
			Type:    "order_failed",
			Payload: map[string]string{"payment_intent_id": paymentIntentID},
		}); err != nil {
			log.Printf("[confirm] ws already closed for session %s: %v", sessionID, err)
		}
		uc.ws.Remove(sessionID.String())
		return nil
	}

	// 2. cancel the asynq expiry task so it doesn't fire after successful payment
	taskID := fmt.Sprintf("%s:%s:%s", sessionID, paymentIntentID, domain.TypeReservationExpired)
	if err := uc.inspector.DeleteTask("default", taskID); err != nil {
		// task may have already fired or doesn't exist — log but don't fail
		log.Printf("[confirm] could not delete asynq task %s: %v", taskID, err)
	}

	// 3. notify the open purchase socket
	if err := uc.ws.Notify(sessionID.String(), sockets.WSMessage{
		Type:    "order_confirmed",
		Payload: map[string]string{"payment_intent_id": paymentIntentID},
	}); err != nil {
		log.Printf("[confirm] ws already closed for session %s: %v", sessionID, err)
	}

	uc.ws.Remove(sessionID.String())

	items, err := uc.purchases.GetTicketIDsByPaymentIntent(ctx, paymentIntentID)
	if err != nil {
		log.Printf("[confirm] failed to fetch ticket ids for %s: %v", paymentIntentID, err)
	} else if len(items) > 0 {
		grants := make([]domain.TicketGrant, 0, len(items))
		for _, item := range items {
			grants = append(grants, domain.TicketGrant{
				TicketID: item.TicketID,
				UserID:   item.UserID,
			})
		}

		var task *asynq.Task
		task, err = domain.NewGrantTicketPermissionsTask(grants, paymentIntentID)
		if err != nil {
			log.Printf("[confirm] failed to create grant permissions task: %v", err)
		} else if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
			log.Printf("[confirm] failed to enqueue grant permissions task: %v", err)
		}
	}

	return nil
}
