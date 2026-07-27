package models

import (
	"time"

	"github.com/google/uuid"
)

type Edition struct {
	ID      uuid.UUID `json:"id"`
	EventID uuid.UUID `json:"event_id"`

	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Tagline     *string `json:"tagline"`
	Description *string `json:"description"`
	IsDraft     bool    `json:"is_draft"`

	RegistrationOpensAt *time.Time `json:"registration_opens_at"`
	StartsAt            time.Time  `json:"starts_at"`
	EndsAt              time.Time  `json:"ends_at"`

	LocationName    *string `json:"location_name"`
	LocationAddress *string `json:"location_description"`

	LogoURL      *string `json:"logo_url"`
	BannerURL    *string `json:"banner_url"`
	ContactEmail *string `json:"contact_email"`

	CreatedBy uuid.UUID  `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateEditionRequest struct {
	Name     string    `json:"name"      validate:"required,min=2"`
	Slug     string    `json:"slug"      validate:"required,min=2"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

func (r CreateEditionRequest) ToInput(eventID uuid.UUID) CreateEditionInput {
	return CreateEditionInput{
		EventID:  eventID,
		Name:     r.Name,
		Slug:     r.Slug,
		StartsAt: r.StartsAt,
		EndsAt:   r.EndsAt,
	}
}

type CreateEditionInput struct {
	EventID  uuid.UUID `json:"event_id"`
	Name     string
	Slug     string
	StartsAt time.Time
	EndsAt   time.Time
}

type PatchEditionRequest struct {
	Name                string     `json:"name"                  validate:"required,min=2"`
	Slug                string     `json:"slug"                  validate:"required,min=2"`
	Tagline             *string    `json:"tagline"`
	Description         *string    `json:"description"`
	RegistrationOpensAt *time.Time `json:"registration_opens_at"`
	StartsAt            time.Time  `json:"starts_at"`
	EndsAt              time.Time  `json:"ends_at"`
	LocationName        *string    `json:"location_name"`
	LocationAddress     *string    `json:"location_description"`
	LogoURL             *string    `json:"logo_url"`
	BannerURL           *string    `json:"banner_url"`
	ContactEmail        *string    `json:"contact_email"`
}

func (r PatchEditionRequest) ToInput(editionID uuid.UUID) PatchEditionInput {
	return PatchEditionInput{
		EditionID:           editionID,
		Name:                r.Name,
		Slug:                r.Slug,
		Tagline:             r.Tagline,
		Description:         r.Description,
		RegistrationOpensAt: r.RegistrationOpensAt,
		StartsAt:            r.StartsAt,
		EndsAt:              r.EndsAt,
		LocationName:        r.LocationName,
		LocationAddress:     r.LocationAddress,
		LogoURL:             r.LogoURL,
		BannerURL:           r.BannerURL,
		ContactEmail:        r.ContactEmail,
	}
}

type PatchEditionInput struct {
	EditionID           uuid.UUID
	Name                string
	Slug                string
	Tagline             *string
	Description         *string
	RegistrationOpensAt *time.Time
	StartsAt            time.Time
	EndsAt              time.Time
	LocationName        *string
	LocationAddress     *string
	LogoURL             *string
	BannerURL           *string
	ContactEmail        *string
}
