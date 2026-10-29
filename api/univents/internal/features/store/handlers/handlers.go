// Package store serves the storefront's live-stock SSE surface (split 6):
// the raw `GET /editions/{edition_id}/store/stream` route (mounted by the
// app router, outside the strict handler — streaming must bypass the
// fun/validate envelope machinery).
package handlers

import (
	"context"
	"net/http"

	"univents/internal/features/store/service"
	"univents/internal/openapi"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

// ListEditionStoreStock is the public `GET /editions/{edition_id}/store/stock`
// route: every purchasable item's current stock position — the REST twin of
// the SSE snapshot payload, so the front parses both with the same shape
// and can graduate to the live stream later without changing anything.
// Unknown editions are NOT_FOUND.
func (h *Handlers) ListEditionStoreStock(ctx context.Context, req openapi.ListEditionStoreStockRequestObject) (openapi.ListEditionStoreStockResponseObject, error) {
	items, err := h.ops.Stock(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	data := make([]openapi.StoreStockItem, 0, len(items))
	for _, it := range items {
		data = append(data, openapi.StoreStockItem{
			Id:       it.ID,
			ItemType: openapi.PurchaseItemType(it.ItemType),
			Stock:    it.Stock,
		})
	}
	return openapi.ListEditionStoreStock200JSONResponse{Items: data}, nil
}

// ServeStoreStream is the raw `GET /editions/{edition_id}/store/stream`
// route (public, D10): the storefront's live stock. Unknown editions get a
// plain 404; known ones stream the snapshot then deltas (hijacked
// close-delimited SSE — the harness WriteTimeout does not apply).
func (h *Handlers) ServeStoreStream(w http.ResponseWriter, r *http.Request) {
	editionID, err := uuid.Parse(chi.URLParam(r, "edition_id"))
	if err != nil {
		http.Error(w, "invalid edition_id", http.StatusBadRequest)
		return
	}
	h.ops.ServeStream(w, r, editionID)
}
