package purchases

import (
	"context"
	"log"
	"net/http"
	"time"
	"univents/internal/platform/telemetry"
	"univents/internal/shared/contracts"
	"univents/internal/shared/sockets"
	"univents/internal/shared/validation"

	"git.trieoh.com/TrieOH/Payssage-SDK-Go"
	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler struct {
	Registry *sockets.Registry
	commands *CommandService
	queries  *QueryService
}

func NewHandler(
	commands *CommandService,
	queries *QueryService,
	registry *sockets.Registry,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
		Registry: registry,
	}
}

func Routes(
	r *chi.Mux,
	h *Handler,
	jwt func(http.Handler) http.Handler,
) {
	r.Post("/webhooks/payments", h.WebhookHandler)
	r.Get("/events/{event_id}/editions/{edition_id}/products/purchase", h.Purchase) // WS upgrade
	r.Route("/purchases", func(r chi.Router) {
		r.Use(jwt)
		r.Get("/", h.ListUserPurchases)
		r.Get("/{purchase_id}/items", h.ListPurchaseItems)
	})
}

// ListUserPurchases godoc
// @Summary List user purchases
// @Description Returns all purchases belonging to the authenticated user
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {object} object
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /purchases [get]
func (handler *Handler) ListUserPurchases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	out, err := handler.queries.ListUserPurchases(ctx)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
}

// ListPurchaseItems godoc
// @Summary List items of a purchase
// @Description Returns all line items belonging to a purchase owned by the authenticated user
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param purchase_id path string true "Purchase ID"
// @Success 200 {object} object
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /purchases/{purchase_id}/items [get]
func (handler *Handler) ListPurchaseItems(w http.ResponseWriter, r *http.Request) {
	purchaseID, rs := validation.GetUUID(r, "purchase_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()

	out, err := handler.queries.ListPurchaseItems(ctx, purchaseID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
}

// Purchase godoc
// @Summary Start a purchase process
// @Description Opens a WebSocket connection to reserve products and initiate payment.
// @Description
// @Description Flow:
// @Description 1. Client sends buy_request {items: [{product_id, quantity}]}
// @Description 2. Server reserves what it can and responds with one of:
// @Description    - reservation_failed: nothing could be reserved (all out of stock)
// @Description    - partial_reservation: some items unavailable, client must respond with confirm_partial or cancel within 60s
// @Description    - reservation_confirmed: all items reserved, proceed to payment
// @Description 3. Client sends submit_payment {card_token, payment_method_id, installments}
// @Description 4. Server responds with one of:
// @Description    - payment_processing: payment submitted, waiting for webhook
// @Description    - payment_failed: payment was rejected
// @Description    - payment_pending: webhook taking too long, poll GET /purchases instead
// @Description    - order_confirmed: payment succeeded (pushed via webhook)
// @Description    - order_failed: payment succeeded but order could not be fulfilled (pushed via webhook)
// @Tags products
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Success 101 {object} object "WebSocket upgrade"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/purchase [get]
func (handler *Handler) Purchase(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		fun.Unauthorized("missing token").Send(w)
		return
	}

	secret := viper.GetString("WS_JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenStr, &contracts.WSClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		fun.Unauthorized("invalid token").Send(w)
		return
	}

	claims, ok := token.Claims.(*contracts.WSClaims)
	if !ok {
		fun.Unauthorized("invalid token claims").Send(w)
		return
	}

	userID := claims.UserID

	upgrader := sockets.MakeUpgrader()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade failed: %v", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			telemetry.Log().Error("failed to close websocket connection", zap.Error(err))
		}
	}(conn)

	spanCtx := trace.SpanFromContext(r.Context()).SpanContext()
	baseCtx := trace.ContextWithSpanContext(context.Background(), spanCtx)

	ctx, cancel := context.WithTimeout(baseCtx, contracts.ReservationDuration+91*time.Second)
	defer cancel()

	if err := handler.commands.Purchase(ctx, conn, editionID, userID); err != nil {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: err.Error()})
		return
	}
}

// WebhookHandler godoc
// @Summary Payment webhook receiver
// @Description Receives normalized payment events from TriePayments. Verifies webhook signature before processing.
// @Description
// @Description Handled events:
// @Description - payment.succeeded: confirms purchase, releases reservation, grants ticket permissions
// @Description - payment.failed: cancels purchase, notifies open WebSocket session if present
// @Description - payment.cancelled: same as payment.failed
// @Description
// @Description Always ACKs with 200 immediately — processing happens asynchronously.
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} object "ACK"
// @Failure 400 {object} contracts.ErrorResponse "Invalid webhook signature"
// @Router /webhooks/payments [post]
func (handler *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := payssage.VerifyWebhookSignature(r, viper.GetString("PAYSSAGE_WEBHOOK_SECRET"))
	if err != nil {
		log.Printf("[webhook] invalid signature: %v", err)
		fun.BadRequest("invalid signature").Send(w)
		return
	}
	log.Printf("[webhook] received event=%s intent=%s", payload.Event, payload.IntentID)

	// ACK immediately — processing is async
	fun.OK().Send(w)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		switch payload.Event {
		case payssage.EventPaymentSucceeded:
			if err := handler.commands.ConfirmPayment(ctx, payload); err != nil {
				log.Printf("[webhook] failed to confirm payment for intent %s: %v", payload.IntentID, err)
			}

		case payssage.EventPaymentFailed, payssage.EventPaymentCancelled:
			if err := handler.commands.CancelPayment(ctx, payload); err != nil {
				log.Printf("[webhook] failed to cancel payment for intent %s: %v", payload.IntentID, err)
			}

		default:
			log.Printf("[webhook] unhandled event type %s for intent %s", payload.Event, payload.IntentID)
		}
	}()
}
