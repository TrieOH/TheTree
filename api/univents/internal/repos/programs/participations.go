package programs

import (
	"context"
	"errors"

	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateParticipation registers a program participation at checkout
// (split 7), attached to the required ticket item's registration — one of
// the purchase's materialized rows (D4). The self-service feature registers
// through the same seam.
func (repo *Repo) CreateParticipation(ctx context.Context, toCreate *models.ProgramParticipation) (*models.ProgramParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.CreateParticipation")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreateProgramParticipation(ctx, sqlc.CreateProgramParticipationParams{
		EditionID:      toCreate.EditionID,
		OccurrenceID:   toCreate.OccurrenceID,
		RegistrationID: toCreate.RegistrationID,
		Status:         string(toCreate.Status),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapParticipation(result)), nil
}

// UpdateParticipationStatus flips a participation (registered→cancelled/
// attended/no_show) — the webhook receiver (split 4) and the expiry worker
// (split 7) write side.
func (repo *Repo) UpdateParticipationStatus(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.UpdateParticipationStatus")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdateProgramParticipationStatus(ctx, sqlc.UpdateProgramParticipationStatusParams{
		ID:     id,
		Status: string(status),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapParticipation(result)), nil
}

// UpsertAttended is the checkpoint check-in write: one atomic statement
// that creates the attended participation on first scan or flips an
// existing live row to attended (idempotent re-scan). The partial unique
// index is the conflict target, so a cancelled row is untouched and the
// insert creates a fresh row (append-only ledger).
func (repo *Repo) UpsertAttended(ctx context.Context, editionID, occurrenceID, registrationID uuid.UUID) (*models.ProgramParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.UpsertAttended")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpsertParticipationAttended(ctx, sqlc.UpsertParticipationAttendedParams{
		EditionID:      editionID,
		OccurrenceID:   occurrenceID,
		RegistrationID: registrationID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapParticipation(result)), nil
}

// GetParticipationByID returns one participation by id.
func (repo *Repo) GetParticipationByID(ctx context.Context, id uuid.UUID) (*models.ProgramParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.GetParticipationByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetProgramParticipationByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapParticipation(result)), nil
}

// GetActiveByOccurrenceAndRegistration returns the registration's live
// participation in an occurrence, or NOT_FOUND when there is none.
func (repo *Repo) GetActiveByOccurrenceAndRegistration(ctx context.Context, occurrenceID, registrationID uuid.UUID) (*models.ProgramParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.GetActiveByOccurrenceAndRegistration")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetActiveProgramParticipationByOccurrenceAndRegistration(ctx, sqlc.GetActiveProgramParticipationByOccurrenceAndRegistrationParams{
		OccurrenceID:   occurrenceID,
		RegistrationID: registrationID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapParticipation(result)), nil
}

// CountActiveByOccurrence is the occupancy count (capacity check inside the
// register tx).
func (repo *Repo) CountActiveByOccurrence(ctx context.Context, occurrenceID uuid.UUID) (int64, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.CountActiveByOccurrence")
	defer span.End()
	return database.Queries(ctx, repo.q).CountActiveProgramParticipationsByOccurrence(ctx, occurrenceID)
}

// UpdateParticipationStatusIfRegistered flips only 'registered' rows — the
// de-registration guard. A nil result (0 rows updated) means the spot is
// locked (attended/no_show) or a racing state change; the caller maps it to
// 409. Mirrors the purchases repo's guarded update (pgx.ErrNoRows on a
// 0-row UPDATE...RETURNING is a guard miss, not a missing row).
func (repo *Repo) UpdateParticipationStatusIfRegistered(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.UpdateParticipationStatusIfRegistered")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdateProgramParticipationStatusIfRegistered(ctx, sqlc.UpdateProgramParticipationStatusIfRegisteredParams{
		ID:     id,
		Status: string(status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//nolint:nilnil // guard missed — the caller maps it to 409
			return nil, nil
		}
		return nil, repo.dbe(err)
	}
	return new(mapParticipation(result)), nil
}

// UpdateParticipationStatusIfNotCancelled flips any non-cancelled row — the
// mark-attended guard (idempotent attended→attended, allows the
// no_show→attended staff correction). A nil result (0 rows updated) means
// cancelled; the caller maps it to 409.
func (repo *Repo) UpdateParticipationStatusIfNotCancelled(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.UpdateParticipationStatusIfNotCancelled")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdateProgramParticipationStatusIfNotCancelled(ctx, sqlc.UpdateProgramParticipationStatusIfNotCancelledParams{
		ID:     id,
		Status: string(status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//nolint:nilnil // guard missed — the caller maps it to 409
			return nil, nil
		}
		return nil, repo.dbe(err)
	}
	return new(mapParticipation(result)), nil
}

// ListActiveByEditionAndRegistration is the "my activities" read: the
// caller's live participations joined with program + occurrence.
func (repo *Repo) ListActiveByEditionAndRegistration(ctx context.Context, editionID, registrationID uuid.UUID) ([]models.MyParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.ListActiveByEditionAndRegistration")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListActiveProgramParticipationsByEditionAndRegistration(ctx, sqlc.ListActiveProgramParticipationsByEditionAndRegistrationParams{
		EditionID:      editionID,
		RegistrationID: registrationID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapMyParticipation), nil
}

// ListByOccurrence is the staff attendance surface: live participations
// with the attendee identity from their registration.
func (repo *Repo) ListByOccurrence(ctx context.Context, occurrenceID uuid.UUID) ([]models.ProgramParticipant, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.ListByOccurrence")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListProgramParticipationsByOccurrence(ctx, occurrenceID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapParticipant), nil
}
