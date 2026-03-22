package domain

import (
	"encoding/json"
	"fmt"
	"time"
	"univents/internal/shared/errx"
	"univents/internal/shared/validation"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type InvalidProduct struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Reason    string    `json:"reason"`
	Requested int       `json:"requested"`
	Reserved  int       `json:"reserved"`
}

type CartItem struct {
	ProductID    uuid.UUID `json:"product_id"`
	Quantity     int       `json:"quantity"`
	HasInventory bool      `json:"-"`
}

type InventoryUpdate struct {
	ProductID          uuid.UUID `json:"product_id"`
	InventoryRemaining int       `json:"inventory_remaining"`
}

type ReservationOutcome struct {
	Reserved         []CartItem        `json:"reserved"`
	Unavailable      []InvalidProduct  `json:"unavailable"`
	InventoryUpdates []InventoryUpdate `json:"inventory_updates"`
}

type ProductType string

const (
	ProductTypeMerchandise ProductType = "merchandise"
	ProductTypeTicket      ProductType = "ticket"
	ProductTypeToken       ProductType = "token"
	ProductTypeBundle      ProductType = "bundle"
)

type ProductStatus string

const (
	ProductStatusDraft       ProductStatus = "draft"
	ProductStatusAvailable   ProductStatus = "available"
	ProductStatusSoldOut     ProductStatus = "sold_out"
	ProductStatusUnavailable ProductStatus = "unavailable"
)

type Product struct {
	ID                 uuid.UUID     `json:"id"`
	ScopeID            uuid.UUID     `json:"scope_id"`
	EditionID          uuid.UUID     `json:"edition_id"`
	Name               string        `json:"name"`
	Description        *string       `json:"description"`
	Type               ProductType   `json:"type"`
	TicketID           *uuid.UUID    `json:"ticket_id"`
	PriceCents         int           `json:"price_cents"`
	Status             ProductStatus `json:"status"`
	AvailableFrom      *time.Time    `json:"available_from"`
	AvailableUntil     *time.Time    `json:"available_until"`
	HasInventory       bool          `json:"has_inventory"`
	InventoryQuantity  int           `json:"inventory_quantity"`
	InventoryRemaining int           `json:"inventory_remaining"`
	CreatedBy          uuid.UUID     `json:"created_by"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	DeletedAt          *time.Time    `json:"deleted_at"`
}

type CreateProductSpec struct {
	EditionID          uuid.UUID   `json:"edition_id"`
	Name               string      `json:"name"`
	Description        *string     `json:"description"`
	Type               ProductType `json:"type"`
	TicketID           *uuid.UUID  `json:"ticket_id"`
	PriceCents         int         `json:"price_cents"`
	AvailableFrom      *time.Time  `json:"available_from"`
	AvailableUntil     *time.Time  `json:"available_until"`
	HasInventory       bool        `json:"has_inventory"`
	InventoryQuantity  int         `json:"inventory_quantity"`
	InventoryRemaining int         `json:"inventory_remaining"`
}

func NewProduct(creatorID uuid.UUID, spec CreateProductSpec) (*Product, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("product").SetMessage("error generating uuid").SetCause(err)
	}

	p := &Product{
		ID:                 id,
		EditionID:          spec.EditionID,
		Name:               spec.Name,
		Description:        spec.Description,
		Type:               spec.Type,
		TicketID:           spec.TicketID,
		PriceCents:         spec.PriceCents,
		Status:             ProductStatusDraft,
		AvailableFrom:      spec.AvailableFrom,
		AvailableUntil:     spec.AvailableUntil,
		HasInventory:       spec.HasInventory,
		InventoryQuantity:  spec.InventoryQuantity,
		InventoryRemaining: spec.InventoryQuantity,
		CreatedBy:          creatorID,
	}

	if err := p.validate(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Product) AddScope(scopeID uuid.UUID) {
	p.ScopeID = scopeID
}

func (p *Product) validate() error {
	return validation.Run(
		validation.RequireUUID("product", "edition_id", p.EditionID),
		validation.RequireUUID("product", "created_by", p.CreatedBy),
		validation.RequireString("product", "name", p.Name),
		validation.RequireString("product", "type", string(p.Type)),
		validation.Assert("product", p.PriceCents >= 0, "invalid price amount"),
		validation.AssertIf("product",
			func() bool { return p.Type == ProductTypeTicket },
			func() bool { return p.TicketID != nil && *p.TicketID != uuid.Nil },
			"if the product is a ticket, ticket_id is required",
		),
		validation.AssertIf("product",
			func() bool { return p.AvailableFrom != nil && p.AvailableUntil != nil },
			func() bool { return p.AvailableFrom.Before(*p.AvailableUntil) },
			"available_from must be before available_until",
		),
		validation.AssertIf("product",
			func() bool { return p.HasInventory },
			func() bool { return p.InventoryQuantity > 0 },
			"product must have at least 1 item available",
		),
		validation.AssertIf("product",
			func() bool { return p.HasInventory },
			func() bool { return p.InventoryQuantity == p.InventoryRemaining },
			"remaining quantity must be equal to starting quantity on creation",
		),
	)
}

const (
	TypeReservationExpired = "reservation:expired"
	ReservationDuration    = 10 * time.Minute
)

type ReservationExpiredPayload struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	EditionID uuid.UUID `json:"edition_id"`
}

func NewReservationExpiredTask(sessionID, userID, editionID uuid.UUID, expiresAt time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(ReservationExpiredPayload{
		SessionID: sessionID,
		UserID:    userID,
		EditionID: editionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal reservation expired payload: %w", err)
	}

	return asynq.NewTask(TypeReservationExpired, payload,
		asynq.TaskID(fmt.Sprintf("%s:%s", sessionID, TypeReservationExpired)),
		asynq.ProcessAt(expiresAt),
		asynq.Unique(time.Hour),
	), nil
}
