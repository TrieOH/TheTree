package domain

import (
	"encoding/json"
	"time"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

type EventStatus string

const (
	StatusDraft        EventStatus = "draft"
	StatusActive       EventStatus = "active"
	StatusArchived     EventStatus = "archived"
	StatusDiscontinued EventStatus = "discontinued"
)

type Event struct {
	ID             uuid.UUID        `json:"id"`
	OwnerID        *uuid.UUID       `json:"owner_id"`
	OrganizationID *uuid.UUID       `json:"organization_id"`
	GoauthScopeID  uuid.UUID        `json:"goauth_scope_id"`
	Name           string           `json:"name"`
	Acronym        *string          `json:"acronym"`
	Slug           string           `json:"slug"`
	Tagline        *string          `json:"tagline"`
	Description    *string          `json:"description"`
	IsSeries       bool             `json:"is_series"`
	EditionsCount  int              `json:"editions_count"`
	LogoUrl        *string          `json:"logo_url"`
	BannerUrl      *string          `json:"banner_url"`
	HasGallery     bool             `json:"has_gallery"`
	GalleryUrls    []string         `json:"gallery_urls"`
	ContactEmail   *string          `json:"contact_email"`
	SocialLinks    *json.RawMessage `json:"social_links"`
	Status         EventStatus      `json:"status"`
	CreatedBy      uuid.UUID        `json:"created_by"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      *time.Time       `json:"deleted_at"`
}

type CreateEventSpec struct {
	OrganizationID *uuid.UUID
	Name           string
	Acronym        *string
	Slug           string
	Tagline        *string
	Description    *string
	IsSeries       bool
	LogoUrl        *string
	ContactEmail   *string
}

type PatchEventSpec struct {
	ID           uuid.UUID
	Name         string
	Acronym      *string
	Slug         string
	Tagline      *string
	Description  *string
	IsSeries     bool
	LogoUrl      *string
	BannerUrl    *string
	HasGallery   bool
	ContactEmail *string
	SocialLinks  *json.RawMessage
}

func NewEvent(creatorID uuid.UUID, ownerID *uuid.UUID, spec CreateEventSpec) (*Event, error) {
	eventUUID, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("event").SetMessage("error generating uuid").SetCause(err)
	}

	now := time.Now()
	return &Event{
		ID:             eventUUID,
		OwnerID:        ownerID,
		OrganizationID: spec.OrganizationID,
		GoauthScopeID:  uuid.Nil,
		Name:           spec.Name,
		Acronym:        spec.Acronym,
		Slug:           spec.Slug,
		Tagline:        spec.Tagline,
		Description:    spec.Description,
		IsSeries:       spec.IsSeries,
		EditionsCount:  0,
		LogoUrl:        spec.LogoUrl,
		BannerUrl:      nil,
		HasGallery:     false,
		GalleryUrls:    nil,
		ContactEmail:   spec.ContactEmail,
		SocialLinks:    nil,
		Status:         StatusDraft,
		CreatedBy:      creatorID,
		CreatedAt:      now,
		UpdatedAt:      now,
		DeletedAt:      nil,
	}, nil
}

func (e *Event) AddScope(scopeID uuid.UUID) {
	e.GoauthScopeID = scopeID
}

type EventAuditAction string

const (
	EventAuditActionCreated              EventAuditAction = "created"
	EventAuditActionEdited               EventAuditAction = "edited"
	EventAuditActionActivated            EventAuditAction = "activated"
	EventAuditActionArchived             EventAuditAction = "archived"
	EventAuditActionDiscontinued         EventAuditAction = "discontinued"
	EventAuditActionDeleted              EventAuditAction = "deleted"
	EventAuditActionRestored             EventAuditAction = "restored"
	EventAuditActionLogoUpdated          EventAuditAction = "logo_updated"
	EventAuditActionBannerUpdated        EventAuditAction = "banner_updated"
	EventAuditActionGalleryUpdated       EventAuditAction = "gallery_updated"
	EventAuditActionNameChanged          EventAuditAction = "name_changed"
	EventAuditActionAcronymChanged       EventAuditAction = "acronym_changed"
	EventAuditActionSlugChanged          EventAuditAction = "slug_changed"
	EventAuditActionTaglineChanged       EventAuditAction = "tagline_changed"
	EventAuditActionDescriptionChanged   EventAuditAction = "description_changed"
	EventAuditActionIsSeriesChanged      EventAuditAction = "is_series_changed"
	EventAuditActionHasGalleryChanged    EventAuditAction = "has_gallery_changed"
	EventAuditActionContactUpdated       EventAuditAction = "contact_updated"
	EventAuditActionSocialLinksUpdated   EventAuditAction = "social_links_updated"
	EventAuditActionEditionAdded         EventAuditAction = "edition_added"
	EventAuditActionEditionRemoved       EventAuditAction = "edition_removed"
	EventAuditActionScopeChanged         EventAuditAction = "scope_changed"
	EventAuditActionOwnershipTransferred EventAuditAction = "ownership_transferred"
)
