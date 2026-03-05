package dtos

import (
	"time"
	"univents/internal/commerce/domain"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	EditionScopeID    uuid.UUID          `json:"edition_scope_id"`
	Name              string             `json:"name" validate:"required,min=3"`
	Description       *string            `json:"description"`
	Type              domain.ProductType `json:"type"`
	TicketID          *uuid.UUID         `json:"ticket_id"`
	PriceCents        int                `json:"price_cents" validate:"gte=0"`
	AvailableFrom     *time.Time         `json:"available_from"`
	AvailableUntil    *time.Time         `json:"available_until"`
	HasInventory      bool               `json:"has_inventory"`
	InventoryQuantity int                `json:"inventory_quantity"`
}
