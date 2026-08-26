package models

import (
	"time"

	"github.com/google/uuid"
)

type RegistrationStatus string

const (
	RegistrationStatusPending   RegistrationStatus = "pending"
	RegistrationStatusConfirmed RegistrationStatus = "confirmed"
	RegistrationStatusCancelled RegistrationStatus = "cancelled"
	RegistrationStatusExpired   RegistrationStatus = "expired"
)

// Registration is one person's enrollment in an edition via a ticket type.
// The checkout feature owns the write side; badges reads confirmed
// registrations to emit participant badges. AttendeeUserID is nil for
// gifted tickets to recipients without an IdentityX account yet (email-only
// gifts) — the recipient claims the ticket after creating an account.
type Registration struct {
	ID               uuid.UUID          `json:"id"`
	EditionID        uuid.UUID          `json:"edition_id"`
	TicketTypeID     uuid.UUID          `json:"ticket_type_id"`
	PurchaserID      uuid.UUID          `json:"purchaser_id"`
	AttendeeUserID   *uuid.UUID         `json:"attendee_user_id"`
	AttendeeEmail    string             `json:"attendee_email"`
	AttendeeName     string             `json:"attendee_name"`
	Status           RegistrationStatus `json:"status"`
	StatusReason     *string            `json:"status_reason"`
	PayssageIntentID *uuid.UUID         `json:"payssage_intent_id"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        *time.Time         `json:"updated_at"`
	DeletedAt        *time.Time         `json:"deleted_at"`
}

// MyTicket is the caller's held ticket for an edition: their own
// registration (pending or confirmed) plus its ticket type. Powers the
// "what do I hold" read the front uses to show upgrade options on the more
// expensive ticket types.
type MyTicket struct {
	RegistrationID uuid.UUID
	TicketType     TicketType
	Status         RegistrationStatus
}
