package dto

import (
	"time"
	"univents/internal/core/domain"
)

type CreateEditionRequest struct {
	Type                 domain.EditionType `json:"type"`
	EditionName          string             `json:"edition_name" validate:"required,min=3,max=256"`
	Tagline              *string            `json:"tagline" validate:"omitempty,max=512"`
	Description          *string            `json:"description" validate:"omitempty,max=8000"`
	RegistrationOpensAt  *time.Time         `json:"registration_opens_at"`
	RegistrationClosesAt *time.Time         `json:"registration_closes_at"`
	StartsAt             time.Time          `json:"starts_at"`
	EndsAt               time.Time          `json:"ends_at"`
	Timezone             string             `json:"timezone"`
	LocationName         string             `json:"location_name"`
	LocationAddress      string             `json:"location_address"`
	LogoUrl              *string            `json:"logo_url" validate:"omitempty,url"`
	BannerUrl            *string            `json:"banner_url" validate:"omitempty,url"`
	ContactEmail         *string            `json:"contact_email" validate:"omitempty,email"`
	ContactPhone         *string            `json:"contact_phone"`
	OrganizerName        *string            `json:"organizer_name"`
}
