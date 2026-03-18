package queries

import (
	"context"
	"univents/internal/commerce/domain"

	"github.com/google/uuid"
)

func (uc *QueryService) StreamInventory(ctx context.Context, editionID uuid.UUID) (<-chan []domain.InventoryUpdate, error) {
	return uc.inventory.Subscribe(ctx, editionID)
}
