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

type CreateBadgeTemplateRequest struct {
	TicketTypeID *uuid.UUID      `json:"ticket_type_id"`
	Name         string          `json:"name" validate:"required,max=256"`
	DesignData   json.RawMessage `json:"design_data" validate:"required"`
}

type CreateBadgeTemplateInput struct {
	EditionID    uuid.UUID
	TicketTypeID *uuid.UUID
	Name         string
	DesignData   json.RawMessage
}

func (req *CreateBadgeTemplateRequest) ToInput(editionID uuid.UUID) CreateBadgeTemplateInput {
	return CreateBadgeTemplateInput{
		EditionID:    editionID,
		TicketTypeID: req.TicketTypeID,
		Name:         req.Name,
		DesignData:   req.DesignData,
	}
}
