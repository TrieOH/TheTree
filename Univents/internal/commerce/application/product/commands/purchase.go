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

func (uc *CommandService) Purchase(ctx context.Context, conn *websocket.Conn, editionID, userID uuid.UUID) error {
	var firstMsg sockets.WSMessage
	if err := conn.ReadJSON(&firstMsg); err != nil {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "invalid message"})
		return nil
	}

	switch firstMsg.Type {
	case "resume_session":
		payloadBytes, err := json.Marshal(firstMsg.Payload)
		if err != nil {
			return err
		}
		var resumeReq dtos.ResumeSessionPayload
		if err = json.Unmarshal(payloadBytes, &resumeReq); err != nil {
			_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "invalid resume_session payload"})
			return nil
		}

		session, err := uc.sessions.Load(ctx, userID, resumeReq.SessionID)
		if err != nil || session == nil || time.Now().After(session.ExpiresAt) {
			_ = conn.WriteJSON(sockets.WSMessage{Type: "session_expired"})
			return nil
		}

		_ = conn.WriteJSON(sockets.WSMessage{
			Type: "reservation_confirmed",
			Payload: dtos.ReservationConfirmedPayload{
				SessionID:     session.SessionID,
				ExpiresAt:     session.ExpiresAt,
				ReservedItems: session.Reserved,
				TotalCents:    session.TotalCents,
			},
		})

		payReq, err := uc.submitPayment(ctx, conn)
		if err != nil || payReq == nil {
			return err
		}

		intent, err := uc.checkout(ctx, conn, session, payReq)
		if err != nil {
			return err
		}

		if err := uc.recordPurchase(ctx, conn, recordPurchaseInput{session: session, intent: intent}); err != nil {
			return err
		}

		return uc.waitForPayment(conn, session)

	case "buy_request":
		payloadBytes, err := json.Marshal(firstMsg.Payload)
		if err != nil {
			return err
		}
		var req dtos.BuyRequest
		if err = json.Unmarshal(payloadBytes, &req); err != nil {
			_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "invalid buy_request payload"})
			return nil
		}

		sessionID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		expiresAt := time.Now().Add(domain.ReservationDuration)

		productMap, err := uc.fetchAndValidateStage(ctx, conn, &req)
		if err != nil {
			return err
		}

		session, err := uc.reserveItemsStage(ctx, conn, reserveItemsInput{
			userID:     userID,
			sessionID:  sessionID,
			editionID:  editionID,
			items:      req.Items,
			expiresAt:  expiresAt,
			productMap: productMap,
		})
		if err != nil || session == nil {
			return err
		}

		payReq, err := uc.submitPayment(ctx, conn)
		if err != nil || payReq == nil {
			return err
		}

		intent, err := uc.checkout(ctx, conn, session, payReq)
		if err != nil {
			return err
		}

		if err := uc.recordPurchase(ctx, conn, recordPurchaseInput{session: session, intent: intent}); err != nil {
			return err
		}

		return uc.waitForPayment(conn, session)

	default:
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "expected buy_request or resume_session"})
		return nil
	}
}

func (uc *CommandService) fetchAndValidateStage(ctx context.Context, conn *websocket.Conn, req *dtos.BuyRequest) (map[uuid.UUID]domain.Product, error) {
	ids := make([]uuid.UUID, len(req.Items))
	for i, item := range req.Items {
		ids[i] = item.ProductID
	}

	toBuy, err := uc.products.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	if len(toBuy) != len(ids) {
		_ = conn.WriteJSON(sockets.WSMessage{
			Type: "purchase_failed",
			Payload: map[string]any{
				"reason":      "one or more products do not exist",
				"product_ids": ids,
			},
		})
		return nil, errors.New("close socket")
	}

	productMap := make(map[uuid.UUID]domain.Product, len(toBuy))
	for _, p := range toBuy {
		productMap[p.ID] = p
	}

	invalid := make([]domain.InvalidProduct, 0)
	for _, p := range toBuy {
		if p.Status != domain.ProductStatusAvailable && p.Status != domain.ProductStatusSoldOut {
			var reason string
			switch p.Status {
			case domain.ProductStatusDraft:
				reason = "product is not yet available"
			case domain.ProductStatusUnavailable:
				reason = "product is unavailable"
			default:
				reason = "product is in invalid state"
			}
			invalid = append(invalid, domain.InvalidProduct{
				ProductID: p.ID,
				Name:      p.Name,
				Reason:    reason,
			})
		}
	}

	if len(invalid) > 0 {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "purchase_failed", Payload: map[string]any{"invalid_products": invalid}})
		return nil, errors.New("close socket")
	}

	for i, item := range req.Items {
		if p, ok := productMap[item.ProductID]; ok {
			req.Items[i].HasInventory = p.HasInventory
		}
	}

	return productMap, nil
}

type reserveItemsInput struct {
	userID     uuid.UUID
	sessionID  uuid.UUID
	editionID  uuid.UUID
	items      []domain.CartItem
	expiresAt  time.Time
	productMap map[uuid.UUID]domain.Product
}

func (uc *CommandService) reserveItemsStage(ctx context.Context, conn *websocket.Conn, in reserveItemsInput) (*domain.PurchaseSession, error) {
	unreserveAndCleanup := func() {
		updates, uErr := uc.products.UnreserveItems(ctx, in.sessionID)
		if uErr != nil {
			telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
		}
		if len(updates) > 0 {
			_ = uc.inventory.Publish(ctx, in.editionID, updates)
		}
		_ = uc.sessions.Delete(ctx, in.userID, in.sessionID)
	}

	outcome, err := uc.products.ReserveItems(ctx, in.sessionID, in.items, in.expiresAt)
	if err != nil {
		unreserveAndCleanup()
		return nil, err
	}

	if len(outcome.InventoryUpdates) > 0 {
		_ = uc.inventory.Publish(ctx, in.editionID, outcome.InventoryUpdates)
	}

	for i, inv := range outcome.Unavailable {
		if p, ok := in.productMap[inv.ProductID]; ok {
			outcome.Unavailable[i].Name = p.Name
		}
	}

	if len(outcome.Reserved) == 0 {
		unreserveAndCleanup()
		_ = conn.WriteJSON(sockets.WSMessage{
			Type:    "reservation_failed",
			Payload: map[string]any{"unavailable": outcome.Unavailable},
		})
		return nil, nil
	}

	// build reservedDetails + total once, reused in partial_reservation and reservation_confirmed
	reservedDetails := make([]domain.ReservedItem, 0, len(outcome.Reserved))
	var total int
	for _, item := range outcome.Reserved {
		p := in.productMap[item.ProductID]
		total += p.PriceCents * item.Quantity
		reservedDetails = append(reservedDetails, domain.ReservedItem{
			ProductID:   p.ID,
			Name:        p.Name,
			Quantity:    item.Quantity,
			PriceCents:  p.PriceCents,
			ProductType: p.Type,
			TicketID:    p.TicketID,
		})
	}

	session := domain.PurchaseSession{
		SessionID:  in.sessionID,
		UserID:     in.userID,
		EditionID:  in.editionID,
		ExpiresAt:  in.expiresAt,
		Reserved:   reservedDetails,
		TotalCents: total,
	}

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
			unreserveAndCleanup()
			return nil, err
		}

		var confirmMsg sockets.WSMessage
		if err = conn.ReadJSON(&confirmMsg); err != nil {
			unreserveAndCleanup()
			return nil, err
		}

		if err = conn.SetReadDeadline(time.Time{}); err != nil {
			unreserveAndCleanup()
			return nil, err
		}

		switch confirmMsg.Type {
		case "confirm_partial":
			// proceed with outcome.Reserved only
		case "cancel":
			unreserveAndCleanup()
			_ = conn.WriteJSON(sockets.WSMessage{Type: "reservation_cancelled"})
			return nil, errors.New("close socket")
		default:
			unreserveAndCleanup()
			_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "expected confirm_partial or cancel"})
			return nil, errors.New("close socket")
		}
	}

	session.Stage = domain.StageAwaitingPayment

	if err := uc.sessions.Save(ctx, session); err != nil {
		return nil, err
	}

	task, err := domain.NewReservationExpiredTask(session.SessionID, session.UserID, session.EditionID, session.ExpiresAt)
	if err != nil {
		unreserveAndCleanup()
		return nil, err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		unreserveAndCleanup()
		return nil, err
	}

	_ = conn.WriteJSON(sockets.WSMessage{
		Type: "reservation_confirmed",
		Payload: dtos.ReservationConfirmedPayload{
			SessionID:     session.SessionID,
			ExpiresAt:     session.ExpiresAt,
			ReservedItems: session.Reserved,
			TotalCents:    session.TotalCents,
		},
	})

	return &session, nil
}

func (uc *CommandService) submitPayment(ctx context.Context, conn *websocket.Conn) (*dtos.SubmitPaymentPayload, error) {
	var payMsg sockets.WSMessage
	if err := conn.ReadJSON(&payMsg); err != nil {
		return nil, errors.New("close socket")
	}

	if payMsg.Type != "submit_payment" {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "expected submit_payment message type"})
		return nil, errors.New("close socket")
	}

	payloadBytes, err := json.Marshal(payMsg.Payload)
	if err != nil {
		return nil, err
	}

	var payReq dtos.SubmitPaymentPayload
	if err = json.Unmarshal(payloadBytes, &payReq); err != nil {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "invalid submit_payment payload"})
		return nil, errors.New("close socket")
	}

	return &payReq, nil
}

func (uc *CommandService) checkout(ctx context.Context, conn *websocket.Conn, session *domain.PurchaseSession, payReq *dtos.SubmitPaymentPayload) (*paymentsSDK.Intent, error) {
	edition, err := uc.editions.GetByID(ctx, session.EditionID)
	if err != nil {
		return nil, err
	}

	telemetry.Log().Info("Before Initiate",
		zap.Int("total", session.TotalCents),
		zap.String("currency", "BRL"),
		zap.String("provider", viper.GetString("TRIEPAYMENTS_PROVIDER")),
		zap.Any("metadata", json.RawMessage(`{"session_id": "`+session.SessionID.String()+`"}`)),
		zap.String("payment_method_id", payReq.PaymentMethodID),
		zap.Int("installments", payReq.Installments),
		zap.String("card_token", payReq.CardToken),
		zap.String("payment_method_type", payReq.PaymentMethodType),
		zap.String("seller_credential_id", edition.TriePaymentsCredentialID.String()),
		zap.String("payer_email", payReq.PayerEmail),
	)

	unreserveAndCleanup := func() {
		updates, uErr := uc.products.UnreserveItems(ctx, session.SessionID)
		if uErr != nil {
			telemetry.Log().Debug("Unreserve failed", zap.Error(uErr))
		}
		if len(updates) > 0 {
			_ = uc.inventory.Publish(ctx, session.EditionID, updates)
		}
		_ = uc.sessions.Delete(ctx, session.UserID, session.SessionID)
	}

	intent, err := uc.payments.InitiateCheckout(ctx, paymentsSDK.InitiateCheckoutRequest{
		Amount:             int64(session.TotalCents),
		Currency:           "BRL",
		Provider:           viper.GetString("TRIEPAYMENTS_PROVIDER"),
		Metadata:           json.RawMessage(`{"session_id": "` + session.SessionID.String() + `"}`),
		PaymentMethodID:    payReq.PaymentMethodID,
		Installments:       payReq.Installments,
		CardToken:          payReq.CardToken,
		PaymentMethodType:  payReq.PaymentMethodType,
		SellerCredentialID: edition.TriePaymentsCredentialID.String(),
		PayerEmail:         payReq.PayerEmail,
	})
	if err != nil {
		unreserveAndCleanup()
		return nil, err
	}

	if intent.MercadoPagoData.PixQRCode != "" {
		if err := uc.sessions.Delete(ctx, session.UserID, session.SessionID); err != nil {
			telemetry.Log().Debug("Failed to delete session after pix checkout", zap.Error(err))
		}
		_ = conn.WriteJSON(sockets.WSMessage{
			Type: "pix_created",
			Payload: map[string]any{
				"qr_code":        intent.MercadoPagoData.PixQRCode,
				"qr_code_base64": intent.MercadoPagoData.PixQRCodeB64,
			},
		})
		return intent, nil
	}

	chargedIntent, err := uc.payments.Charge(ctx, intent.ID, paymentsSDK.ChargeRequest{
		SellerCredentialID: edition.TriePaymentsCredentialID.String(),
	})
	if err != nil {
		unreserveAndCleanup()
		return nil, err
	}

	if err := uc.sessions.Delete(ctx, session.UserID, session.SessionID); err != nil {
		telemetry.Log().Debug("Failed to delete session after checkout", zap.Error(err))
	}

	_ = conn.WriteJSON(sockets.WSMessage{Type: "payment_processing"})
	return chargedIntent, nil
}

type recordPurchaseInput struct {
	session *domain.PurchaseSession
	intent  *paymentsSDK.Intent
}

// FIXME Use outbox pattern to avoid user losing purchase on TX failure

func (uc *CommandService) recordPurchase(ctx context.Context, conn *websocket.Conn, in recordPurchaseInput) error {
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

func (uc *CommandService) waitForPayment(conn *websocket.Conn, session *domain.PurchaseSession) error {
	uc.ws.Register(session.SessionID.String(), conn)

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
