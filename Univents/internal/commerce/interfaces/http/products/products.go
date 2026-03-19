package products

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"univents/internal/commerce/application/product/commands"
	"univents/internal/commerce/application/product/queries"
	"univents/internal/commerce/domain"
	"univents/internal/commerce/interfaces/http/dtos"
	"univents/internal/plataform/telemetry"
	"univents/internal/shared/sockets"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler struct {
	Registry *sockets.Registry
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewProductsHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
	registry *sockets.Registry,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
		Registry: registry,
	}
}

// Create godoc
// @Summary Create a new product
// @Description Creates a new product for an edition.
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param request body dtos.CreateProductRequest true "Product creation request"
// @Success 201 {object} object "Product created successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateProductRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := domain.CreateProductSpec{
		EditionID:          editionID,
		Name:               req.Name,
		Description:        req.Description,
		Type:               req.Type,
		TicketID:           req.TicketID,
		PriceCents:         req.PriceCents,
		AvailableFrom:      req.AvailableFrom,
		AvailableUntil:     req.AvailableUntil,
		HasInventory:       req.HasInventory,
		InventoryQuantity:  req.InventoryQuantity,
		InventoryRemaining: req.InventoryQuantity,
	}

	ctx := r.Context()
	out, err := handler.commands.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(out).Send(w)
}

// Publish godoc
// @Summary publishes a product
// @Description Publishes a product making it publicly available.
// @Tags activities
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} object "Product published successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/{product_id}/publish [post]
func (handler *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	productID, rs := validation.GetUUID(r, "product_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.Publish(ctx, productID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}

// List godoc
// @Summary List all edition products
// @Description List all publicly available products of the edition
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Success 201 {object} object
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products [get]
func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.List(ctx, editionID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
}

// ListAdmin godoc
// @Summary List all edition products
// @Description if user has permission products:read list all edition products
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Success 201 {object} object
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/admin [get]
func (handler *Handler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	out, err := handler.queries.AdminList(ctx, editionID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
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
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /purchases [get]
func (handler *Handler) ListUserPurchases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	out, err := handler.queries.ListUserPurchases(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
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
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 403 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(out).Send(w)
}

var upgrader = websocket.Upgrader{}

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
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/purchase [get]
func (handler *Handler) Purchase(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		resp.Unauthorized("missing token").Send(w)
		return
	}

	secret := viper.GetString("WS_JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenStr, &domain.WSClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		resp.Unauthorized("invalid token").Send(w)
		return
	}

	claims, ok := token.Claims.(*domain.WSClaims)
	if !ok {
		resp.Unauthorized("invalid token claims").Send(w)
		return
	}

	userID := claims.UserID
	userEmail := claims.Email

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

	var req dtos.BuyRequest
	if err := conn.ReadJSON(&req); err != nil {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "invalid request payload"})
		return
	}

	spanCtx := trace.SpanFromContext(r.Context()).SpanContext()
	baseCtx := trace.ContextWithSpanContext(context.Background(), spanCtx)

	ctx, cancel := context.WithTimeout(baseCtx, domain.ReservationDuration+91*time.Second)
	defer cancel()

	if err := handler.commands.Purchase(ctx, conn, req, editionID, userID, userEmail); err != nil {
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
// @Failure 400 {object} swag.ErrorResponse "Invalid webhook signature"
// @Router /webhooks/payments [post]
func (handler *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := paymentsSDK.VerifyWebhookSignature(r, viper.GetString("TRIEPAYMENTS_WEBHOOK_SECRET"))
	if err != nil {
		log.Printf("[webhook] invalid signature: %v", err)
		resp.BadRequest("invalid signature").Send(w)
		return
	}
	log.Printf("[webhook] received event=%s intent=%s", payload.Event, payload.IntentID)

	// ACK immediately — processing is async
	resp.OK().Send(w)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		switch payload.Event {
		case paymentsSDK.EventPaymentSucceeded:
			if err := handler.commands.ConfirmPayment(ctx, payload.IntentID); err != nil {
				log.Printf("[webhook] failed to confirm payment for intent %s: %v", payload.IntentID, err)
			}

		case paymentsSDK.EventPaymentFailed, paymentsSDK.EventPaymentCancelled:
			if err := handler.commands.CancelPayment(ctx, payload.IntentID); err != nil {
				log.Printf("[webhook] failed to cancel payment for intent %s: %v", payload.IntentID, err)
			}

		default:
			log.Printf("[webhook] unhandled event type %s for intent %s", payload.Event, payload.IntentID)
		}
	}()
}

// StreamInventory godoc
// @Summary Stream inventory updates for an edition store
// @Description Opens a Server-Sent Events stream that pushes inventory_update events whenever
// @Description product stock changes due to reservations, cancellations, or expiries.
// @Description
// @Description Event format:
// @Description   event: inventory_update
// @Description   data: [{"product_id": "...", "inventory_remaining": 3}, ...]
// @Description
// @Description The stream stays open until the client disconnects.
// @Tags products
// @Produce text/event-stream
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Success 200 {object} object "SSE stream"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/inventory/stream [get]
func (handler *Handler) StreamInventory(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		resp.InternalServerError("streaming not supported").Send(w)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

	updates, err := handler.queries.StreamInventory(ctx, editionID)
	if err != nil {
		resp.InternalServerError("failed to subscribe to inventory stream").Send(w)
		return
	}

	// send initial ping so client knows stream is alive
	_, _ = fmt.Fprintf(w, ": ping\n\n")
	flusher.Flush()

	keepalive := time.NewTicker(29 * time.Second)
	defer keepalive.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-keepalive.C:
			_, _ = fmt.Fprintf(w, ": ping\n\n")
			flusher.Flush()
		case batch, ok := <-updates:
			if !ok {
				return
			}
			payload, err := json.Marshal(batch)
			if err != nil {
				continue
			}
			_, _ = fmt.Fprintf(w, "event: inventory_update\ndata: %s\n\n", payload)
			flusher.Flush()
		}
	}
}
