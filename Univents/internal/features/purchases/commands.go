package purchases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"univents/internal/platform/database"
	"univents/internal/platform/telemetry"
	"univents/internal/shared/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"
	"univents/internal/shared/sockets"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	editions  ports.EditionsRepository
	products  ports.ProductsRepository
	purchases ports.PurchaseRepository
	payments  *paymentsSDK.Client
	sessions  ports.PurchaseSessionStore
	ws        *sockets.Registry
	inventory ports.InventoryPublisher
	minio     *minio.Client
	asynq     *asynq.Client
	inspector *asynq.Inspector
	tracer    trace.Tracer
	az        *authzed.Client
	tx        database.TxRunner
}

func NewCommandService(
	editions ports.EditionsRepository,
	products ports.ProductsRepository,
	purchases ports.PurchaseRepository,
	payments *paymentsSDK.Client,
	session ports.PurchaseSessionStore,
	ws *sockets.Registry,
	inventory ports.InventoryPublisher,
	minio *minio.Client,
	asynq *asynq.Client,
	inspector *asynq.Inspector,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		editions:  editions,
		products:  products,
		purchases: purchases,
		payments:  payments,
		sessions:  session,
		ws:        ws,
		inventory: inventory,
		minio:     minio,
		asynq:     asynq,
		inspector: inspector,
		tracer:    tracer,
		az:        az,
		tx:        tx,
	}
}

type ResumeSessionPayload struct {
	SessionID uuid.UUID `json:"session_id"`
}

type BuyRequest struct {
	Items []contracts.CartItem `json:"items"`
}

type ReservationConfirmedPayload struct {
	SessionID     uuid.UUID                `json:"session_id"`
	ExpiresAt     time.Time                `json:"expires_at"`
	ReservedItems []contracts.ReservedItem `json:"reserved_items"`
	TotalCents    int                      `json:"total_cents"`
}

type SubmitPaymentPayload struct {
	CardToken            string `json:"card_token"`
	PaymentMethodID      string `json:"payment_method_id"`
	PaymentMethodType    string `json:"payment_method_type"`
	Installments         int    `json:"installments"`
	PayerEmail           string `json:"payer_email"`
	IdentificationNumber string `json:"identification_number"`
	IdentificationType   string `json:"identification_type"`
}

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
		var resumeReq ResumeSessionPayload
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
			Payload: ReservationConfirmedPayload{
				SessionID:     session.SessionID,
				ExpiresAt:     session.ExpiresAt,
				ReservedItems: session.Reserved,
				TotalCents:    session.TotalCents,
			},
		})

		payReq, err := uc.submitPayment(conn)
		if err != nil {
			return err
		}
		if payReq == nil {
			// user canceled before submitting payment
			updates, _ := uc.products.UnreserveItems(ctx, session.SessionID)
			if len(updates) > 0 {
				_ = uc.inventory.Publish(ctx, session.EditionID, updates)
			}
			_ = uc.sessions.Delete(ctx, session.UserID, session.SessionID)
			_ = conn.WriteJSON(sockets.WSMessage{Type: "purchase_cancelled"})
			return nil
		}

		intent, isPix, err := uc.checkout(ctx, conn, session, payReq)
		if err != nil {
			return err
		}

		return uc.waitForPayment(ctx, conn, session, intent, isPix)

	case "buy_request":
		payloadBytes, err := json.Marshal(firstMsg.Payload)
		if err != nil {
			return err
		}
		var req BuyRequest
		if err = json.Unmarshal(payloadBytes, &req); err != nil {
			_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "invalid buy_request payload"})
			return nil
		}

		sessionID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		expiresAt := time.Now().Add(contracts.ReservationDuration)

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

		payReq, err := uc.submitPayment(conn)
		if err != nil {
			return err
		}
		if payReq == nil {
			// user canceled before submitting payment
			updates, _ := uc.products.UnreserveItems(ctx, session.SessionID)
			if len(updates) > 0 {
				_ = uc.inventory.Publish(ctx, session.EditionID, updates)
			}
			_ = uc.sessions.Delete(ctx, session.UserID, session.SessionID)
			_ = conn.WriteJSON(sockets.WSMessage{Type: "purchase_cancelled"})
			return nil
		}

		intent, isPix, err := uc.checkout(ctx, conn, session, payReq)
		if err != nil {
			return err
		}

		return uc.waitForPayment(ctx, conn, session, intent, isPix)

	default:
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "expected buy_request or resume_session"})
		return nil
	}
}

func (uc *CommandService) fetchAndValidateStage(ctx context.Context, conn *websocket.Conn, req *BuyRequest) (map[uuid.UUID]contracts.Product, error) {
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

	productMap := make(map[uuid.UUID]contracts.Product, len(toBuy))
	for _, p := range toBuy {
		productMap[p.ID] = p
	}

	invalid := make([]contracts.InvalidProduct, 0)
	for _, p := range toBuy {
		if p.Status != contracts.ProductStatusAvailable && p.Status != contracts.ProductStatusSoldOut {
			var reason string
			switch p.Status {
			case contracts.ProductStatusDraft:
				reason = "product is not yet available"
			case contracts.ProductStatusUnavailable:
				reason = "product is unavailable"
			default:
				reason = "product is in invalid state"
			}
			invalid = append(invalid, contracts.InvalidProduct{
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
	items      []contracts.CartItem
	expiresAt  time.Time
	productMap map[uuid.UUID]contracts.Product
}

func (uc *CommandService) reserveItemsStage(ctx context.Context, conn *websocket.Conn, in reserveItemsInput) (*contracts.CheckoutSession, error) {
	unreserveAndCleanup := func() {
		updates, uErr := uc.products.UnreserveItems(ctx, in.sessionID)
		if uErr != nil {
			telemetry.Log().Info("Unreserve failed", zap.Error(uErr))
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
	reservedDetails := make([]contracts.ReservedItem, 0, len(outcome.Reserved))
	var total int
	for _, item := range outcome.Reserved {
		p := in.productMap[item.ProductID]
		total += p.PriceCents * item.Quantity
		reservedDetails = append(reservedDetails, contracts.ReservedItem{
			ProductID:   p.ID,
			Name:        p.Name,
			Quantity:    item.Quantity,
			PriceCents:  p.PriceCents,
			ProductType: p.Type,
			TicketID:    p.TicketID,
		})
	}

	session := contracts.CheckoutSession{
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

	session.Stage = contracts.StageAwaitingPayment

	if err := uc.sessions.Save(ctx, session); err != nil {
		return nil, err
	}

	task, err := contracts.NewReservationExpiredTask(session.SessionID, session.UserID, session.EditionID, session.ExpiresAt)
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
		Payload: ReservationConfirmedPayload{
			SessionID:     session.SessionID,
			ExpiresAt:     session.ExpiresAt,
			ReservedItems: session.Reserved,
			TotalCents:    session.TotalCents,
		},
	})

	return &session, nil
}

func (uc *CommandService) submitPayment(conn *websocket.Conn) (*SubmitPaymentPayload, error) {
	var payMsg sockets.WSMessage
	if err := conn.ReadJSON(&payMsg); err != nil {
		return nil, errors.New("close socket")
	}

	switch payMsg.Type {
	case "submit_payment":
		payloadBytes, err := json.Marshal(payMsg.Payload)
		if err != nil {
			return nil, err
		}
		var payReq SubmitPaymentPayload
		if err = json.Unmarshal(payloadBytes, &payReq); err != nil {
			_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "invalid submit_payment payload"})
			return nil, errors.New("close socket")
		}
		return &payReq, nil

	case "cancel_purchase":
		return nil, nil // nil, nil signals caller to run cleanup

	default:
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "expected submit_payment or cancel_purchase"})
		return nil, errors.New("close socket")
	}
}

func (uc *CommandService) checkout(ctx context.Context, conn *websocket.Conn, session *contracts.CheckoutSession, payReq *SubmitPaymentPayload) (*paymentsSDK.Intent, bool, error) {
	edition, err := uc.editions.GetByID(ctx, session.EditionID)
	if err != nil {
		return nil, false, err
	}

	unreserveAndCleanup := func() {
		updates, uErr := uc.products.UnreserveItems(ctx, session.SessionID)
		if uErr != nil {
			telemetry.Log().Info("Unreserve failed", zap.Error(uErr))
		}
		if len(updates) > 0 {
			_ = uc.inventory.Publish(ctx, session.EditionID, updates)
		}
		_ = uc.sessions.Delete(ctx, session.UserID, session.SessionID)
	}

	if edition.TriePaymentsProvider == nil {
		telemetry.Log().Error("edition not set up to receive payments, should block store")
		unreserveAndCleanup()
		return nil, true, errors.New("edition not set up to receive payments")
	}

	telemetry.Log().Info("Before Initiate",
		zap.Int("total", session.TotalCents),
		zap.String("currency", "BRL"),
		zap.String("provider", *edition.TriePaymentsProvider),
		zap.Any("metadata", json.RawMessage(`{"session_id": "`+session.SessionID.String()+`", "user_id": "`+session.UserID.String()+`"}`)),
		zap.String("payment_method_id", payReq.PaymentMethodID),
		zap.Int("installments", payReq.Installments),
		zap.String("card_token", payReq.CardToken),
		zap.String("payment_method_type", payReq.PaymentMethodType),
		zap.String("seller_credential_id", edition.TriePaymentsCredentialID.String()),
		zap.String("payer_email", payReq.PayerEmail),
		zap.String("identification_number", payReq.IdentificationNumber),
		zap.String("identification_type", payReq.IdentificationType),
	)

	intent, err := uc.payments.InitiateCheckout(ctx, paymentsSDK.InitiateCheckoutRequest{
		Amount:               int64(session.TotalCents),
		Currency:             "BRL",
		Provider:             *edition.TriePaymentsProvider,
		Metadata:             json.RawMessage(`{"session_id": "` + session.SessionID.String() + `", "user_id": "` + session.UserID.String() + `"}`),
		PaymentMethodID:      payReq.PaymentMethodID,
		Installments:         payReq.Installments,
		CardToken:            payReq.CardToken,
		PaymentMethodType:    payReq.PaymentMethodType,
		SellerCredentialID:   edition.TriePaymentsCredentialID.String(),
		PayerEmail:           payReq.PayerEmail,
		IdentificationNumber: payReq.IdentificationNumber,
		IdentificationType:   payReq.IdentificationType,
	})
	if err != nil {
		unreserveAndCleanup()
		return nil, false, err
	}

	if intent.MercadoPagoData.PixQRCode != "" {
		_ = conn.WriteJSON(sockets.WSMessage{
			Type: "pix_created",
			Payload: map[string]any{
				"qr_code":        intent.MercadoPagoData.PixQRCode,
				"qr_code_base64": intent.MercadoPagoData.PixQRCodeB64,
			},
		})
		return intent, true, nil
	}

	_ = conn.WriteJSON(sockets.WSMessage{Type: "payment_processing"})
	return intent, false, nil
}

func (uc *CommandService) cancelPixRequest(ctx context.Context, conn *websocket.Conn, session *contracts.CheckoutSession, intent *paymentsSDK.Intent) error {
	edition, err := uc.editions.GetByID(ctx, session.EditionID)
	if err != nil {
		telemetry.Log().Info("Failed to fetch edition for pix cancel", zap.Error(err))
	} else {
		telemetry.Log().Info("Trying to cancel pix payment", zap.String("intent_id", intent.ID))
		if _, err := uc.payments.CancelPixIntent(ctx, intent.ID, paymentsSDK.CancelPixRequest{
			Provider:           intent.Provider,
			SellerCredentialID: edition.TriePaymentsCredentialID.String(),
		}); err != nil {
			telemetry.Log().Info("Failed to cancel pix intent", zap.Error(err))
		}
	}

	updates, err := uc.products.UnreserveItems(ctx, session.SessionID)
	if err != nil {
		telemetry.Log().Info("Unreserve failed on pix cancel", zap.Error(err))
	}
	if len(updates) > 0 {
		_ = uc.inventory.Publish(ctx, session.EditionID, updates)
	}

	if err := uc.sessions.Delete(ctx, session.UserID, session.SessionID); err != nil {
		telemetry.Log().Info("Failed to delete session on pix cancel", zap.Error(err))
	}

	uc.ws.Remove(session.SessionID.String())

	_ = conn.WriteJSON(sockets.WSMessage{Type: "purchase_cancelled"})
	return nil
}

func (uc *CommandService) waitForPayment(ctx context.Context, conn *websocket.Conn, session *contracts.CheckoutSession, intent *paymentsSDK.Intent, isPix bool) error {
	uc.ws.Register(session.SessionID.String(), conn)

	paymentTimeout := time.Until(session.ExpiresAt)
	if !isPix {
		paymentTimeout = 30 * time.Second
	}

	timer := time.After(paymentTimeout)

	connClosed := make(chan struct{})
	webhookMsg := make(chan sockets.WSMessage, 1)
	cancelMsg := make(chan struct{}, 1)

	uc.ws.RegisterCallback(session.SessionID.String(), func(msg sockets.WSMessage) {
		webhookMsg <- msg
	})

	go func() {
		for {
			var msg sockets.WSMessage
			if err := conn.ReadJSON(&msg); err != nil {
				close(connClosed)
				return
			}
			if isPix && msg.Type == "cancel_purchase" {
				cancelMsg <- struct{}{}
				return
			}
		}
	}()

	select {
	case msg := <-webhookMsg:
		_ = conn.WriteJSON(msg)
		return nil
	case <-cancelMsg:
		return uc.cancelPixRequest(ctx, conn, session, intent)
	case <-timer:
		_ = conn.WriteJSON(sockets.WSMessage{
			Type:    "payment_pending",
			Payload: "payment is taking longer than expected, you can close this and check your purchases",
		})
		return nil
	case <-connClosed:
		return nil
	}
}

func (uc *CommandService) ConfirmPayment(ctx context.Context, payload *paymentsSDK.WebhookPayload) error {
	paymentIntentID := payload.IntentID

	var meta struct {
		SessionID uuid.UUID `json:"session_id"`
		UserID    uuid.UUID `json:"user_id"`
	}
	if err := json.Unmarshal(payload.Metadata, &meta); err != nil || meta.SessionID == uuid.Nil || meta.UserID == uuid.Nil {
		return fmt.Errorf("missing session_id or user_id in metadata for intent %s", paymentIntentID)
	}

	session, err := uc.sessions.Load(ctx, meta.UserID, meta.SessionID)
	if err != nil || session == nil {
		return fmt.Errorf("session not found for intent %s", paymentIntentID)
	}

	if err := uc.recordPurchase(ctx, recordPurchaseInput{
		session: session,
		intent:  &paymentsSDK.Intent{ID: paymentIntentID, Amount: int64(session.TotalCents), Provider: payload.Provider},
	}); err != nil {
		return fmt.Errorf("failed to record purchase for intent %s: %w", paymentIntentID, err)
	}

	if err := uc.sessions.Delete(ctx, meta.UserID, meta.SessionID); err != nil {
		log.Printf("[confirm] failed to delete session %s: %v", meta.SessionID, err)
	}

	return uc.finalizeConfirmedPurchase(ctx, paymentIntentID)
}

func (uc *CommandService) finalizeConfirmedPurchase(ctx context.Context, paymentIntentID string) error {
	purchase, err := uc.purchases.GetByPaymentID(ctx, paymentIntentID)
	if err != nil {
		return fmt.Errorf("failed to fetch purchase for intent %s: %w", paymentIntentID, err)
	}

	switch purchase.Status {
	case contracts.PurchaseStatusCompleted:
		log.Printf("[confirm] purchase already completed for intent %s", paymentIntentID)
		return nil
	case contracts.PurchaseStatusCancelled:
		log.Printf("[confirm] purchase cancelled, ignoring success webhook for %s", paymentIntentID)
		return nil
	}

	// 1. confirm purchase + clean reservation in one TX
	if err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		if purchase.SessionID != nil {
			if err := uc.products.DeleteReservation(ctx, *purchase.SessionID); err != nil {
				return err
			}
		}
		if err := uc.purchases.ConfirmPurchase(ctx, paymentIntentID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to confirm purchase tx for intent %s: %w", paymentIntentID, err)
	}

	// 2. cancel expiry task
	if purchase.SessionID != nil {
		taskID := fmt.Sprintf("%s:%s", *purchase.SessionID, contracts.TypeReservationExpired)
		if err := uc.inspector.DeleteTask("default", taskID); err != nil {
			log.Printf("[confirm] could not delete asynq task %s: %v", taskID, err)
		}
	}

	// 3. notify WS if still alive
	if purchase.SessionID != nil {
		sessionID := *purchase.SessionID
		if err := uc.ws.Notify(sessionID.String(), sockets.WSMessage{
			Type:    "payment_confirmed",
			Payload: map[string]string{"purchase_id": purchase.ID.String()},
		}); err != nil {
			log.Printf("[confirm] ws already closed for session %s: %v", sessionID, err)
		}
		uc.ws.Remove(sessionID.String())
	}

	// 4. grant ticket permissions
	items, err := uc.purchases.GetTicketIDsByPaymentIntent(ctx, paymentIntentID)
	if err != nil {
		log.Printf("[confirm] failed to fetch ticket ids for %s: %v", paymentIntentID, err)
		return nil
	}

	if len(items) > 0 {
		grants := make([]contracts.TicketGrant, 0, len(items))
		for _, item := range items {
			grants = append(grants, contracts.TicketGrant{
				TicketID: item.TicketID,
				UserID:   item.UserID,
			})
		}
		task, err := contracts.NewGrantTicketPermissionsTask(grants, paymentIntentID)
		if err != nil {
			log.Printf("[confirm] failed to create grant permissions task: %v", err)
			return nil
		}
		if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
			log.Printf("[confirm] failed to enqueue grant permissions task: %v", err)
		}
	}

	return nil
}

func (uc *CommandService) CancelPayment(ctx context.Context, payload *paymentsSDK.WebhookPayload) error {
	paymentIntentID := payload.IntentID

	var meta struct {
		SessionID uuid.UUID `json:"session_id"`
		UserID    uuid.UUID `json:"user_id"`
	}
	if err := json.Unmarshal(payload.Metadata, &meta); err != nil || meta.SessionID == uuid.Nil || meta.UserID == uuid.Nil {
		return fmt.Errorf("missing session_id or user_id in metadata for intent %s", paymentIntentID)
	}

	if err := uc.sessions.Delete(ctx, meta.UserID, meta.SessionID); err != nil {
		log.Printf("[cancel] failed to delete session %s: %v", meta.SessionID, err)
	}

	if err := uc.ws.Notify(meta.SessionID.String(), sockets.WSMessage{
		Type:    "payment_failed",
		Payload: map[string]string{"payment_intent_id": paymentIntentID},
	}); err != nil {
		log.Printf("[cancel] ws already closed for session %s: %v", meta.SessionID, err)
	}
	uc.ws.Remove(meta.SessionID.String())

	return nil
}

type recordPurchaseInput struct {
	session *contracts.CheckoutSession
	intent  *paymentsSDK.Intent
}

func (uc *CommandService) recordPurchase(ctx context.Context, in recordPurchaseInput) error {
	if err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		pendingPurchase := contracts.NewPurchase(contracts.CreatePurchaseSpec{
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
			if item.ProductType == contracts.ProductTypeTicket {
				for range item.Quantity {
					if _, err = uc.purchases.CreateLineItem(ctx, contracts.LineItem{
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
				if _, err = uc.purchases.CreateLineItem(ctx, contracts.LineItem{
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
