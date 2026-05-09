package products

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"univents/internal/shared/contracts"
	"univents/internal/shared/validation"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func Routes(
	r *chi.Mux,
	h *Handler,
	jwt func(http.Handler) http.Handler,
) {
	r.Route("/events/{event_id}/editions/{edition_id}/products", func(r chi.Router) {
		r.Get("/", h.List)
		r.Use(jwt)
		r.Post("/", h.Create)
		r.Get("/admin", h.ListAdmin)
		r.Get("/inventory/stream", h.StreamInventory) // SSE upgrade
		r.Route("/{product_id}", func(r chi.Router) {
			r.Delete("/", h.Delete)
			r.Post("/publish", h.Publish)
			r.Post("/restore", h.Restore)
			r.Post("/gallery", h.AddGalleryImage)
			r.Delete("/gallery", h.RemoveGalleryImage)
			r.Put("/thumbnail", h.SetThumbnail)
			r.Delete("/thumbnail", h.UnsetThumbnail)
		})
	})
}

type CreateProductRequest struct {
	EditionScopeID    uuid.UUID             `json:"edition_scope_id"`
	Name              string                `json:"name" validate:"required,min=3"`
	Description       *string               `json:"description"`
	Type              contracts.ProductType `json:"type"`
	TicketID          *uuid.UUID            `json:"ticket_id"`
	PriceCents        int                   `json:"price_cents" validate:"gte=0"`
	AvailableFrom     *time.Time            `json:"available_from"`
	AvailableUntil    *time.Time            `json:"available_until"`
	HasInventory      bool                  `json:"has_inventory"`
	InventoryQuantity int                   `json:"inventory_quantity"`
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
// @Param request body CreateProductRequest true "Product creation request"
// @Success 201 {object} object "Product created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := contracts.CreateProductSpec{
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
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(out).Send(w)
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
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
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
		fun.Error(err).Send(w)
		return
	}

	fun.OK().Send(w)
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
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
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
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
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
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
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
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
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
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/inventory/stream [get]
func (handler *Handler) StreamInventory(w http.ResponseWriter, r *http.Request) {
	editionID, rs := validation.GetUUID(r, "edition_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		fun.InternalServerError("streaming not supported").Send(w)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

	updates, err := handler.queries.StreamInventory(ctx, editionID)
	if err != nil {
		fun.InternalServerError("failed to subscribe to inventory stream").Send(w)
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

// Delete godoc
// @Summary Soft delete a product
// @Description Soft deletes a product. Blocked if the product has pending or completed purchases.
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} object "Product deleted successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/{product_id} [delete]
func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	productID, rs := validation.GetUUID(r, "product_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	if err := handler.commands.Delete(ctx, productID); err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Product deleted successfully").Send(w)
}

// Restore godoc
// @Summary Restore a soft deleted product
// @Description Restores a soft deleted product. Only works if the product has not been hard deleted yet.
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} object "Product restored successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/{product_id}/restore [post]
func (handler *Handler) Restore(w http.ResponseWriter, r *http.Request) {
	productID, rs := validation.GetUUID(r, "product_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	if err := handler.commands.Restore(ctx, productID); err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Product restored successfully").Send(w)
}

type ImageURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// AddGalleryImage godoc
// @Summary Add an image to the product gallery
// @Description Adds a MinIO URL to the product's gallery_urls array.
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param product_id path string true "Product ID"
// @Param request body ImageURLRequest true "Image URL"
// @Success 200 {object} contracts.Product "Image added to gallery"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/{product_id}/gallery [post]
func (handler *Handler) AddGalleryImage(w http.ResponseWriter, r *http.Request) {
	productID, rs := validation.GetUUID(r, "product_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.AddGalleryImage(ctx, productID, req.URL)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Image added to gallery").WithData(product).Send(w)
}

// RemoveGalleryImage godoc
// @Summary Remove an image from the product gallery
// @Description Removes a URL from the product's gallery_urls array and deletes the object from MinIO.
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param product_id path string true "Product ID"
// @Param request body ImageURLRequest true "Image URL"
// @Success 200 {object} contracts.Product "Image removed from gallery"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/{product_id}/gallery [delete]
func (handler *Handler) RemoveGalleryImage(w http.ResponseWriter, r *http.Request) {
	productID, rs := validation.GetUUID(r, "product_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.RemoveGalleryImage(ctx, productID, req.URL)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Image removed from gallery").WithData(product).Send(w)
}

// SetThumbnail godoc
// @Summary Set the product thumbnail
// @Description Sets the product thumbnail URL. If the URL is not already in gallery_urls it is added automatically.
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param product_id path string true "Product ID"
// @Param request body ImageURLRequest true "Image URL"
// @Success 200 {object} contracts.Product "Thumbnail set"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/{product_id}/thumbnail [put]
func (handler *Handler) SetThumbnail(w http.ResponseWriter, r *http.Request) {
	productID, rs := validation.GetUUID(r, "product_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req ImageURLRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.SetThumbnail(ctx, productID, req.URL)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Thumbnail set").WithData(product).Send(w)
}

// UnsetThumbnail godoc
// @Summary Unset the product thumbnail
// @Description Clears the product thumbnail. The image remains in gallery_urls.
// @Tags products
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param event_id path string true "Event ID"
// @Param edition_id path string true "Edition ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} contracts.Product "Thumbnail unset"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 403 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /events/{event_id}/editions/{edition_id}/products/{product_id}/thumbnail [delete]
func (handler *Handler) UnsetThumbnail(w http.ResponseWriter, r *http.Request) {
	productID, rs := validation.GetUUID(r, "product_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	product, err := handler.commands.UnsetThumbnail(ctx, productID)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK("Thumbnail unset").WithData(product).Send(w)
}
