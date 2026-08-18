package store

import (
	"context"
	"time"

	"univents/internal/openapi"
)

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
	return openapi.ListEditionStoreStock200JSONResponse{
		Code: 200, Data: &data, Timestamp: time.Now(), Module: module,
	}, nil
}
