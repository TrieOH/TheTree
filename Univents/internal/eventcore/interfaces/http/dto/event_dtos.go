package dto

import (
	"github.com/google/uuid"
)

type CreateEventRequest struct {
	OrganizationID *uuid.UUID `json:"organization_id"`
	Name           string     `json:"name" validate:"required,min=2"`
	Acronym        *string    `json:"acronym"`
	Slug           string     `json:"slug" validate:"required,min=2"`
	Tagline        *string    `json:"tagline"`
	Description    *string    `json:"description"`
	IsSeries       bool       `json:"is_series"`
	LogoUrl        *string    `json:"logo_url"`
	ContactEmail   *string    `json:"contact_email" validate:"required,email"`
}
