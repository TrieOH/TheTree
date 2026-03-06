package products

import (
	"net/http"
	"univents/internal/commerce/application/product/commands"
	"univents/internal/commerce/application/product/queries"
	"univents/internal/commerce/domain"
	"univents/internal/commerce/interfaces/http/dtos"
	"univents/internal/shared/validation"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type Handler struct {
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
		EditionScopeID:     req.EditionScopeID,
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
