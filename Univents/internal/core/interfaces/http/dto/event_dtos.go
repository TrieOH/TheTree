package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	GoAuthEventScopeID uuid.UUID  `json:"go_auth_event_scope_id"`
	OrganizationID     *uuid.UUID `json:"organization_id"`
	Name               string     `json:"name" validate:"required,min=2"`
	Acronym            *string    `json:"acronym"`
	Slug               string     `json:"slug" validate:"required,min=2"`
	Tagline            *string    `json:"tagline"`
	Description        *string    `json:"description"`
	IsSeries           bool       `json:"is_series"`
	LogoUrl            *string    `json:"logo_url"`
	BannerUrl          *string    `json:"banner_url"`
	ContactEmail       *string    `json:"contact_email" validate:"required,email"`
}

type PatchEventRequest struct {
	Name         string           `json:"name" validate:"required,min=3,max=256"`
	Acronym      *string          `json:"acronym" validate:"omitempty,min=2,max=32"`
	Slug         string           `json:"slug" validate:"required,min=2,max=32"`
	Tagline      *string          `json:"tagline" validate:"omitempty,max=512"`
	Description  *string          `json:"description"`
	IsSeries     bool             `json:"is_series"`
	LogoUrl      *string          `json:"logo_url" validate:"omitempty,url"`
	BannerUrl    *string          `json:"banner_url" validate:"omitempty,url"`
	HasGallery   bool             `json:"has_gallery"`
	ContactEmail *string          `json:"contact_email" validate:"omitempty,email"`
	SocialLinks  *json.RawMessage `json:"social_links" validate:"omitempty,json"`
}
