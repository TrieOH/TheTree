package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	EventStatusDraft        Status = "draft"
	EventStatusActive       Status = "active"
	EventStatusArchived     Status = "archived"
	EventStatusDiscontinued Status = "discontinued"
)

type Event struct {
	ID             uuid.UUID        `json:"id"`
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
	Status         Status           `json:"status"`
	CreatedBy      uuid.UUID        `json:"created_by"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      *time.Time       `json:"deleted_at"`
}

type ActorType string

const (
	ActorTypeParticipant ActorType = "participant"
	ActorTypeStaff       ActorType = "staff"
	ActorTypeAdmin       ActorType = "admin"
	ActorTypeOwner       ActorType = "owner"
	ActorTypeSystem      ActorType = "system"
)

type AuditAction string

const (
	EventAuditActionCreated              AuditAction = "created"
	EventAuditActionEdited               AuditAction = "edited"
	EventAuditActionActivated            AuditAction = "activated"
	EventAuditActionArchived             AuditAction = "archived"
	EventAuditActionDiscontinued         AuditAction = "discontinued"
	EventAuditActionDeleted              AuditAction = "deleted"
	EventAuditActionRestored             AuditAction = "restored"
	EventAuditActionLogoUpdated          AuditAction = "logo_updated"
	EventAuditActionBannerUpdated        AuditAction = "banner_updated"
	EventAuditActionGalleryUpdated       AuditAction = "gallery_updated"
	EventAuditActionSlugChanged          AuditAction = "slug_changed"
	EventAuditActionContactUpdated       AuditAction = "contact_updated"
	EventAuditActionSocialLinksUpdated   AuditAction = "social_links_updated"
	EventAuditActionEditionAdded         AuditAction = "edition_added"
	EventAuditActionEditionRemoved       AuditAction = "edition_removed"
	EventAuditActionScopeChanged         AuditAction = "scope_changed"
	EventAuditActionOwnershipTransferred AuditAction = "ownership_transferred"
)

type Audit struct {
	ID         uuid.UUID        `json:"id"`
	EventID    uuid.UUID        `json:"event_id"`
	ActorType  ActorType        `json:"actor_type"`
	ActorID    *uuid.UUID       `json:"actor_id"`
	Action     AuditAction      `json:"action"`
	FromStatus *Status          `json:"from_status"`
	ToStatus   *Status          `json:"to_status"`
	Metadata   *json.RawMessage `json:"metadata"`
	CreatedAt  time.Time        `json:"created_at"`
}
