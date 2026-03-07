package dtos

import (
	"univents/internal/commerce/domain"

	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Description *string `json:"description"`
}

type AddTicketPermissionRequest struct {
	PermissionType domain.PermissionType `json:"permission_type"`
	ActivityID     *uuid.UUID            `json:"activity_id"`
	ProductID      *uuid.UUID            `json:"product_id"`
	CheckpointID   *uuid.UUID            `json:"checkpoint_id"`
}
