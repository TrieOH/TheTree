package models

import (
	"time"

	"github.com/google/uuid"
)

// TermsOfService is a project's current terms document. It is versioned in
// place: an update bumps Version and refreshes EffectiveAt (the 30-day
// notice period), while the acceptance ledger keeps the per-version proof.
// A project without a row has no terms (v0) — nothing to accept, nothing
// to notify.
type TermsOfService struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"project_id"`
	Version     int       `json:"version"`
	Content     string    `json:"content"`
	EffectiveAt time.Time `json:"effective_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TosAcceptance is one actor's recorded agreement to one version of a
// project's terms. ContentHash pins the SHA-256 of the content at
// acceptance time — the controller bears the burden of proving consent
// under LGPD art. 8 §2, so the row is self-contained evidence. Source
// records HOW consent was given.
type TosAcceptance struct {
	ID          uuid.UUID `json:"id"`
	ActorID     uuid.UUID `json:"actor_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	TosVersion  int       `json:"tos_version"`
	ContentHash string    `json:"content_hash"`
	Source      string    `json:"source"` // "clickwrap" | "continued_use"
	AcceptedAt  time.Time `json:"accepted_at"`
}

const (
	TosAcceptanceClickwrap    = "clickwrap"
	TosAcceptanceContinuedUse = "continued_use"
)

// TosAcceptanceWithActor is an acceptance joined with the actor's email —
// the shape the admin audit listing serves.
type TosAcceptanceWithActor struct {
	TosAcceptance

	ActorEmail *string `json:"actor_email"`
}
