package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BadgeTemplate struct {
	ID           uuid.UUID       `json:"id"`
	EditionID    uuid.UUID       `json:"edition_id"`
	TicketTypeID *uuid.UUID      `json:"ticket_type_id"`
	Name         string          `json:"name"`
	DesignData   json.RawMessage `json:"design_data"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    *time.Time      `json:"updated_at"`
	DeletedAt    *time.Time      `json:"deleted_at"`
}

type CreateBadgeTemplateInput struct {
	EditionID    uuid.UUID
	TicketTypeID *uuid.UUID
	Name         string
	DesignData   json.RawMessage
}
