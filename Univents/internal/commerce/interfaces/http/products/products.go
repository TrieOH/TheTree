package products

import (
	"context"
	"log"
	"net/http"
	"univents/internal/commerce/application/product/commands"
	"univents/internal/commerce/application/product/queries"
	"univents/internal/commerce/domain"
	"univents/internal/commerce/interfaces/http/dtos"
	"univents/internal/shared/sockets"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
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
// @Description Opens a WebSocket connection to reserve products and initiate payment. Sends reservation_confirmed with client_secret on success, reservation_failed if items are unavailable.
// @Tags products
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	var req dtos.BuyRequest
	if err := conn.ReadJSON(&req); err != nil {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: "invalid request payload"})
		return
	}

	ctx := r.Context()
	if err := handler.commands.Purchase(ctx, conn, req, editionID); err != nil {
		_ = conn.WriteJSON(sockets.WSMessage{Type: "error", Payload: err.Error()})
		return
	}

	// Block until conn closes (resolved by webhook or asynq task)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Printf("[buy] ws closed: %v", err)
			return
		}
	}
}

// WebhookHandler godoc
// @Summary TrieMint webhook receiver
// @Description Receives normalized payment events from TrieMint
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} object
// @Failure 400 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /webhooks/payments [post]
func (handler *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := paymentsSDK.VerifyWebhookSignature(r, viper.GetString("TRIEPAYMENTS_WEBHOOK_SECRET"))
	if err != nil {
		log.Printf("[webhook] invalid signature: %v", err)
		resp.BadRequest("invalid signature").Send(w)
		return
	}
	log.Printf("[webhook] received event=%s intent=%s", payload.Event, payload.IntentID)

	// ACK immediately
	resp.OK().Send(w)

	go func() {
		ctx := context.Background()
		switch payload.Event {
		case paymentsSDK.EventPaymentSucceeded:
			if err := handler.commands.ConfirmPayment(ctx, payload.IntentID); err != nil {
				log.Printf("[webhook] failed to confirm payment for intent %s: %v", payload.IntentID, err)
			}
		case paymentsSDK.EventPaymentFailed, paymentsSDK.EventPaymentCancelled:
			purchase, err := handler.queries.GetPurchaseByPaymentID(ctx, payload.IntentID)
			if err != nil {
				log.Printf("[webhook] failed to fetch purchase for intent %s: %v", payload.IntentID, err)
				return
			}
			if purchase.SessionID == nil {
				return
			}
			if err := handler.Registry.Notify(purchase.SessionID.String(), sockets.WSMessage{
				Type:    "payment_failed",
				Payload: map[string]string{"payment_intent_id": payload.IntentID},
			}); err != nil {
				log.Printf("[webhook] ws already closed for session %s: %v", purchase.SessionID, err)
			}
			handler.Registry.Remove(purchase.SessionID.String())
		}
	}()
}
