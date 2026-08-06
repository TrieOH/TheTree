package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BadgeTemplateOrigin scopes a badge template. nil means the default scope
// (edition default or per-ticket-type); BadgeTemplateOriginStaff is the
// edition's staff-only design.
type BadgeTemplateOrigin string

const BadgeTemplateOriginStaff BadgeTemplateOrigin = "staff"

type BadgeTemplate struct {
	ID           uuid.UUID            `json:"id"`
	EditionID    uuid.UUID            `json:"edition_id"`
	TicketTypeID *uuid.UUID           `json:"ticket_type_id"`
	Origin       *BadgeTemplateOrigin `json:"origin"`
	Name         string               `json:"name"`
	DesignData   json.RawMessage      `json:"design_data"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    *time.Time           `json:"updated_at"`
	DeletedAt    *time.Time           `json:"deleted_at"`
}

type CreateBadgeTemplateInput struct {
	EditionID    uuid.UUID
	TicketTypeID *uuid.UUID
	Origin       *BadgeTemplateOrigin
	Name         string
	DesignData   json.RawMessage
}

// UpdateBadgeTemplateInput is the PATCH body: nil fields are left unchanged.
// ticket_type_id and origin are immutable (scope is set at creation and
// changed by delete + create) and cannot be patched.
type UpdateBadgeTemplateInput struct {
	TemplateID uuid.UUID
	Name       *string
	DesignData json.RawMessage // nil = unchanged
}

type BadgeEmissionOrigin string

const (
	BadgeEmissionOriginParticipant BadgeEmissionOrigin = "participant"
	BadgeEmissionOriginStaff       BadgeEmissionOrigin = "staff"
)

type BadgeEmissionStatus string

const (
	BadgeEmissionStatusActive  BadgeEmissionStatus = "active"
	BadgeEmissionStatusRevoked BadgeEmissionStatus = "revoked"
)

// BadgeEmission is one badge awarded to a person for an edition. The design is
// never stored on the emission — it is derived at read time from the current
// badge_templates of the edition (see templateFor). RegistrationID anchors
// participant emissions (name, ticket, template); staff emissions keep it nil.
type BadgeEmission struct {
	ID             uuid.UUID           `json:"id"`
	EditionID      uuid.UUID           `json:"edition_id"`
	UserID         uuid.UUID           `json:"user_id"`
	Origin         BadgeEmissionOrigin `json:"origin"`
	RegistrationID *uuid.UUID          `json:"registration_id"`
	Status         BadgeEmissionStatus `json:"status"`
	StatusReason   *string             `json:"status_reason"`
	EmailSentAt    *time.Time          `json:"email_sent_at"`
	EmittedAt      time.Time           `json:"emitted_at"`
	UpdatedAt      *time.Time          `json:"updated_at"`
}

// BadgeEmissionView is a read row joining the emission with its edition, event
// and ticket data needed for rendering. Names are deliberately absent — the
// front derives holder names from IdentityX profiles via getActorProfile.
type BadgeEmissionView struct {
	BadgeEmission

	EditionName  string
	EndsAt       time.Time
	EventName    string
	TicketTypeID *uuid.UUID
	TicketName   *string
}

type BadgeOriginGroup struct {
	Current []BadgeProfileBadge `json:"current"`
	Past    []BadgeProfileBadge `json:"past"`
}

// BadgeProfileGroups is the public profile grouping: two origens, each split
// into current (edition ends_at >= now) and past editions, most current first.
type BadgeProfileGroups struct {
	Attendant BadgeOriginGroup `json:"attendant"`
	Staff     BadgeOriginGroup `json:"staff"`
}

// BadgeProfileBadge is one badge as shown on a profile, self-contained for the
// front renderer (design + variables). Holder names are derived by the front
// from IdentityX profiles (getActorProfile) — never exposed here.
type BadgeProfileBadge struct {
	EmissionID   uuid.UUID           `json:"emission_id"`
	EditionID    uuid.UUID           `json:"edition_id"`
	EditionName  string              `json:"edition_name"`
	EventName    string              `json:"event_name"`
	Origin       BadgeEmissionOrigin `json:"origin"`
	TemplateID   *uuid.UUID          `json:"template_id"`
	TemplateName *string             `json:"template_name"`
	DesignData   json.RawMessage     `json:"design_data"`
	TicketName   *string             `json:"ticket_name"`
	ActionURL    string              `json:"action_url"`
}

// BadgeEditionEmission is an emission as browsed by the edition owner. Holder
// names are derived by the front from IdentityX profiles — not exposed here.
type BadgeEditionEmission struct {
	ID           uuid.UUID           `json:"id"`
	UserID       uuid.UUID           `json:"user_id"`
	Origin       BadgeEmissionOrigin `json:"origin"`
	Status       BadgeEmissionStatus `json:"status"`
	StatusReason *string             `json:"status_reason"`
	TicketName   *string             `json:"ticket_name"`
	TemplateID   *uuid.UUID          `json:"template_id"`
	TemplateName *string             `json:"template_name"`
	EmittedAt    time.Time           `json:"emitted_at"`
	UpdatedAt    *time.Time          `json:"updated_at"`
}

// BadgePrintItem is the per-emission payload the front renders into the print
// PDF: the design plus the variables resolved from Univents' own data. Holder
// names are derived by the front from IdentityX profiles — not exposed here.
type BadgePrintItem struct {
	EmissionID   uuid.UUID           `json:"emission_id"`
	UserID       uuid.UUID           `json:"user_id"`
	Origin       BadgeEmissionOrigin `json:"origin"`
	EventName    string              `json:"event_name"`
	EditionName  string              `json:"edition_name"`
	TicketName   *string             `json:"ticket_name"`
	TemplateID   *uuid.UUID          `json:"template_id"`
	TemplateName *string             `json:"template_name"`
	DesignData   json.RawMessage     `json:"design_data"`
	ActionURL    string              `json:"action_url"`
}
