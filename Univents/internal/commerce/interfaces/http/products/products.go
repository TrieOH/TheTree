package products

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"univents/internal/commerce/application/product/commands"
	"univents/internal/commerce/application/product/queries"
	"univents/internal/commerce/domain"
	"univents/internal/commerce/interfaces/http/dtos"
	"univents/internal/shared/sockets"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/gorilla/websocket"
)

type Handler struct {
	Registry *sockets.Registry
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewProductsHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
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

// ConfirmPayment godoc
// @Summary Payment webhook callback
// @Description Called by the payments service when a payment intent is confirmed
// @Tags products
// @Accept json
// @Produce json
// @Param body body dtos.ConfirmPaymentRequest true "Payment confirmation"
// @Success 200 {object} object
// @Failure 400 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/purchase/confirm [post]
func (handler *Handler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	var req dtos.ConfirmPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.BadRequest("invalid payload").Send(w)
		return
	}

	resp.OK().Send(w)

	ctx := context.Background()
	if err := handler.commands.ConfirmPayment(ctx, req.SessionID, req.PaymentIntentID); err != nil {
		log.Printf("[webhook] failed to confirm payment for session %s: %v", req.SessionID, err)
	}
}
