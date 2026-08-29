package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

// ProgramParticipationRepo is the program_participations surface. The write
// side has two owners: the checkout flow (created registered at checkout,
// flipped by the webhook receiver / expiry worker) and the self-service
// activity registration feature (register / de-register / mark-attended).
// The guarded updates are the concurrency backstops — the row itself holds
// the state machine, and a 0-row update means the caller raced a state
// change (mapped to 409 by the service). Method names mirror the programs
// repo (satisfied by *programs.Repo).
type ProgramParticipationRepo interface {
	CreateParticipation(ctx context.Context, toCreate *models.ProgramParticipation) (*models.ProgramParticipation, error)
	UpdateParticipationStatus(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error)

	// GetParticipationByID returns one participation by id (staff marking
	// resolves the edition for the role check). NOT_FOUND when unknown.
	GetParticipationByID(ctx context.Context, id uuid.UUID) (*models.ProgramParticipation, error)

	// GetActiveByOccurrenceAndRegistration returns the registration's live
	// participation in an occurrence (register pre-check → 409 already
	// registered; de-register lookup → 404 not registered). NOT_FOUND when
	// there is none.
	GetActiveByOccurrenceAndRegistration(ctx context.Context, occurrenceID, registrationID uuid.UUID) (*models.ProgramParticipation, error)

	// CountActiveByOccurrence is the occupancy count (capacity check inside
	// the register tx).
	CountActiveByOccurrence(ctx context.Context, occurrenceID uuid.UUID) (int64, error)

	// UpdateParticipationStatusIfRegistered flips only 'registered' rows —
	// the de-registration guard (an attended/no_show spot is locked). 0
	// rows = locked or a racing state change; the caller maps it to 409.
	UpdateParticipationStatusIfRegistered(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error)

	// UpdateParticipationStatusIfNotCancelled flips any non-cancelled row —
	// the mark-attended guard (idempotent attended→attended, allows the
	// no_show→attended staff correction). 0 rows = cancelled → 409.
	UpdateParticipationStatusIfNotCancelled(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error)

	// UpsertAttended is the checkpoint check-in write (staff): creates the
	// attended participation on first scan or flips an existing live row to
	// attended — one atomic statement, never duplicates (the partial unique
	// index is the conflict target).
	UpsertAttended(ctx context.Context, editionID, occurrenceID, registrationID uuid.UUID) (*models.ProgramParticipation, error)

	// ListActiveByEditionAndRegistration is the "my activities" read: the
	// caller's live participations joined with program + occurrence.
	ListActiveByEditionAndRegistration(ctx context.Context, editionID, registrationID uuid.UUID) ([]models.MyParticipation, error)

	// ListByOccurrence is the staff attendance surface: live participations
	// with the attendee identity from their registration.
	ListByOccurrence(ctx context.Context, occurrenceID uuid.UUID) ([]models.ProgramParticipant, error)
}
