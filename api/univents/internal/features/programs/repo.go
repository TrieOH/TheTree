package programs

import (
	"context"
	"errors"
	"lib/database"
	"univents/errs"
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ProgramParticipationRepo = (*Repo)(nil)

var (
	_ ports.ProgramRepo           = (*Repo)(nil)
	_ ports.ProgramOccurrenceRepo = (*Repo)(nil)
)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrProgramNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapProgram(src sqlc.Program) models.Program {
	return models.Program{
		ID:             src.ID,
		EditionID:      src.EditionID,
		Kind:           models.ProgramKind(src.Kind),
		Name:           src.Name,
		Description:    src.Description,
		MinAccessLevel: src.MinAccessLevel,
		StaffOnly:      src.StaffOnly,
		Price:          &src.Price,
		BannerURL:      src.BannerUrl,
		CreatedAt:      src.CreatedAt,
		UpdatedAt:      src.UpdatedAt,
		DeletedAt:      src.DeletedAt,
	}
}

func mapProgramOccurrence(src sqlc.ProgramOccurrence) models.ProgramOccurrence {
	return models.ProgramOccurrence{
		ID:          src.ID,
		ProgramID:   src.ProgramID,
		EditionID:   src.EditionID,
		StartsAt:    src.StartsAt,
		EndsAt:      src.EndsAt,
		MaxCapacity: src.MaxCapacity,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}

func mapParticipation(src sqlc.ProgramParticipation) models.ProgramParticipation {
	return models.ProgramParticipation{
		ID:             src.ID,
		EditionID:      src.EditionID,
		OccurrenceID:   src.OccurrenceID,
		RegistrationID: src.RegistrationID,
		Status:         models.ProgramParticipationStatus(src.Status),
		CreatedAt:      src.CreatedAt,
		UpdatedAt:      src.UpdatedAt,
	}
}

func mapMyParticipation(src sqlc.ListActiveProgramParticipationsByEditionAndRegistrationRow) models.MyParticipation {
	return models.MyParticipation{
		ID:             src.ParticipationID,
		EditionID:      src.EditionID,
		OccurrenceID:   src.OccurrenceID,
		RegistrationID: src.RegistrationID,
		Status:         models.ProgramParticipationStatus(src.Status),
		CreatedAt:      src.ParticipationCreatedAt,
		UpdatedAt:      src.ParticipationUpdatedAt,
		Program: models.Program{
			ID:             src.ProgramID,
			EditionID:      src.EditionID,
			Kind:           models.ProgramKind(src.ProgramKind),
			Name:           src.ProgramName,
			Description:    src.ProgramDescription,
			MinAccessLevel: src.ProgramMinAccessLevel,
			StaffOnly:      src.ProgramStaffOnly,
			Price:          &src.ProgramPrice,
			BannerURL:      src.ProgramBannerUrl,
		},
		Occurrence: models.ProgramOccurrence{
			ID:          src.OccurrenceID,
			ProgramID:   src.ProgramID,
			EditionID:   src.EditionID,
			StartsAt:    src.OccurrenceStartsAt,
			EndsAt:      src.OccurrenceEndsAt,
			MaxCapacity: src.OccurrenceMaxCapacity,
		},
	}
}

func mapParticipant(src sqlc.ListProgramParticipationsByOccurrenceRow) models.ProgramParticipant {
	return models.ProgramParticipant{
		ID:             src.ParticipationID,
		OccurrenceID:   src.OccurrenceID,
		RegistrationID: src.RegistrationID,
		Status:         models.ProgramParticipationStatus(src.Status),
		CreatedAt:      src.ParticipationCreatedAt,
		AttendeeUserID: src.AttendeeUserID,
		AttendeeEmail:  src.AttendeeEmail,
		AttendeeName:   src.AttendeeName,
	}
}

func priceValue(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func (repo *Repo) Create(ctx context.Context, toCreate *models.Program) (*models.Program, error) {
	result, err := repo.db(ctx).CreateProgram(ctx, sqlc.CreateProgramParams{
		EditionID:      toCreate.EditionID,
		Kind:           string(toCreate.Kind),
		Name:           toCreate.Name,
		Description:    toCreate.Description,
		MinAccessLevel: toCreate.MinAccessLevel,
		StaffOnly:      toCreate.StaffOnly,
		Price:          priceValue(toCreate.Price),
	})
	return database.One(result, err, repo.dbe, mapProgram)
}

func (repo *Repo) CreateOccurrence(ctx context.Context, toCreate *models.ProgramOccurrence) (*models.ProgramOccurrence, error) {
	result, err := repo.db(ctx).CreateProgramOccurrence(ctx, sqlc.CreateProgramOccurrenceParams{
		ProgramID:   toCreate.ProgramID,
		EditionID:   toCreate.EditionID,
		StartsAt:    toCreate.StartsAt,
		EndsAt:      toCreate.EndsAt,
		MaxCapacity: toCreate.MaxCapacity,
	})
	return database.One(result, err, repo.dbe, mapProgramOccurrence)
}

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	result, err := repo.db(ctx).DeleteProgram(ctx, id)
	return database.One(result, err, repo.dbe, mapProgram)
}

func (repo *Repo) DeleteOccurrence(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	result, err := repo.db(ctx).DeleteProgramOccurrence(ctx, id)
	return database.One(result, err, repo.dbe, mapProgramOccurrence)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	result, err := repo.db(ctx).GetProgramByID(ctx, id)
	return database.One(result, err, repo.dbe, mapProgram)
}

// GetByIDForUpdate is the row-lock variant used inside the checkout tx
// (split 7): serializes concurrent checkouts on the same program (price
// read) before availability is checked.
func (repo *Repo) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	result, err := repo.db(ctx).GetProgramByIDForUpdate(ctx, id)
	return database.One(result, err, repo.dbe, mapProgram)
}

func (repo *Repo) GetOccurrenceByID(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	result, err := repo.db(ctx).GetProgramOccurrenceByID(ctx, id)
	return database.One(result, err, repo.dbe, mapProgramOccurrence)
}

// GetOccurrenceByIDForUpdate is the row-lock variant used inside the
// checkout tx (split 7): serializes concurrent checkouts on the same
// occurrence (capacity) before availability is checked.
func (repo *Repo) GetOccurrenceByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	result, err := repo.db(ctx).GetProgramOccurrenceByIDForUpdate(ctx, id)
	return database.One(result, err, repo.dbe, mapProgramOccurrence)
}

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Program, error) {
	results, err := repo.db(ctx).ListProgramsByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapProgram)
}

func (repo *Repo) ListOccurrencesByProgram(ctx context.Context, programID uuid.UUID) ([]models.ProgramOccurrence, error) {
	results, err := repo.db(ctx).ListProgramOccurrencesByProgram(ctx, programID)
	return database.Many(results, err, repo.dbe, mapProgramOccurrence)
}

func (repo *Repo) ListOccurrencesByEdition(ctx context.Context, editionID uuid.UUID) ([]models.ProgramOccurrence, error) {
	results, err := repo.db(ctx).ListProgramOccurrencesByEdition(ctx, editionID)
	return database.Many(results, err, repo.dbe, mapProgramOccurrence)
}

// CreateParticipation registers a program participation at checkout
// (split 7), attached to the required ticket item's registration — one of
// the purchase's materialized rows (D4). The self-service feature registers
// through the same seam.
func (repo *Repo) CreateParticipation(ctx context.Context, toCreate *models.ProgramParticipation) (*models.ProgramParticipation, error) {
	result, err := repo.db(ctx).CreateProgramParticipation(ctx, sqlc.CreateProgramParticipationParams{
		EditionID:      toCreate.EditionID,
		OccurrenceID:   toCreate.OccurrenceID,
		RegistrationID: toCreate.RegistrationID,
		Status:         string(toCreate.Status),
	})
	return database.One(result, err, repo.dbe, mapParticipation)
}

// UpdateParticipationStatus flips a participation (registered→cancelled/
// attended/no_show) — the webhook receiver (split 4) and the expiry worker
// (split 7) write side.
func (repo *Repo) UpdateParticipationStatus(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error) {
	result, err := repo.db(ctx).UpdateProgramParticipationStatus(ctx, sqlc.UpdateProgramParticipationStatusParams{
		ID:     id,
		Status: string(status),
	})
	return database.One(result, err, repo.dbe, mapParticipation)
}

// UpsertAttended is the checkpoint check-in write: one atomic statement
// that creates the attended participation on first scan or flips an
// existing live row to attended (idempotent re-scan). The partial unique
// index is the conflict target, so a cancelled row is untouched and the
// insert creates a fresh row (append-only ledger).
func (repo *Repo) UpsertAttended(ctx context.Context, editionID, occurrenceID, registrationID uuid.UUID) (*models.ProgramParticipation, error) {
	result, err := repo.db(ctx).UpsertParticipationAttended(ctx, sqlc.UpsertParticipationAttendedParams{
		EditionID:      editionID,
		OccurrenceID:   occurrenceID,
		RegistrationID: registrationID,
	})
	return database.One(result, err, repo.dbe, mapParticipation)
}

// GetParticipationByID returns one participation by id.
func (repo *Repo) GetParticipationByID(ctx context.Context, id uuid.UUID) (*models.ProgramParticipation, error) {
	result, err := repo.db(ctx).GetProgramParticipationByID(ctx, id)
	return database.One(result, err, repo.dbe, mapParticipation)
}

// GetActiveByOccurrenceAndRegistration returns the registration's live
// participation in an occurrence, or database.ErrNotFound when there is
// none.
func (repo *Repo) GetActiveByOccurrenceAndRegistration(ctx context.Context, occurrenceID, registrationID uuid.UUID) (*models.ProgramParticipation, error) {
	result, err := repo.db(ctx).GetActiveProgramParticipationByOccurrenceAndRegistration(ctx, sqlc.GetActiveProgramParticipationByOccurrenceAndRegistrationParams{
		OccurrenceID:   occurrenceID,
		RegistrationID: registrationID,
	})
	return database.One(result, err, repo.dbe, mapParticipation)
}

// CountActiveByOccurrence is the occupancy count (capacity check inside the
// register tx).
func (repo *Repo) CountActiveByOccurrence(ctx context.Context, occurrenceID uuid.UUID) (int64, error) {
	count, err := repo.db(ctx).CountActiveProgramParticipationsByOccurrence(ctx, occurrenceID)
	return count, repo.dbe(err)
}

// UpdateParticipationStatusIfRegistered flips only 'registered' rows — the
// de-registration guard. A nil result (0 rows updated) means the spot is
// locked (attended/no_show) or a racing state change; the caller maps it to
// 409. Mirrors the purchases repo's guarded update (pgx.ErrNoRows on a
// 0-row UPDATE...RETURNING is a guard miss, not a missing row).
func (repo *Repo) UpdateParticipationStatusIfRegistered(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error) {
	result, err := repo.db(ctx).UpdateProgramParticipationStatusIfRegistered(ctx, sqlc.UpdateProgramParticipationStatusIfRegisteredParams{
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
	result, err := repo.db(ctx).UpdateProgramParticipationStatusIfNotCancelled(ctx, sqlc.UpdateProgramParticipationStatusIfNotCancelledParams{
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
	rows, err := repo.db(ctx).ListActiveProgramParticipationsByEditionAndRegistration(ctx, sqlc.ListActiveProgramParticipationsByEditionAndRegistrationParams{
		EditionID:      editionID,
		RegistrationID: registrationID,
	})
	return database.Many(rows, err, repo.dbe, mapMyParticipation)
}

// ListByOccurrence is the staff attendance surface: live participations
// with the attendee identity from their registration.
func (repo *Repo) ListByOccurrence(ctx context.Context, occurrenceID uuid.UUID) ([]models.ProgramParticipant, error) {
	rows, err := repo.db(ctx).ListProgramParticipationsByOccurrence(ctx, occurrenceID)
	return database.Many(rows, err, repo.dbe, mapParticipant)
}

func (repo *Repo) Patch(ctx context.Context, id uuid.UUID, program *models.Program) (*models.Program, error) {
	result, err := repo.db(ctx).PatchProgram(ctx, sqlc.PatchProgramParams{
		Kind:           string(program.Kind),
		Name:           program.Name,
		Description:    program.Description,
		MinAccessLevel: program.MinAccessLevel,
		StaffOnly:      program.StaffOnly,
		Price:          priceValue(program.Price),
		BannerUrl:      program.BannerURL,
		ID:             id,
	})
	return database.One(result, err, repo.dbe, mapProgram)
}

func (repo *Repo) PatchOccurrence(ctx context.Context, id uuid.UUID, occurrence *models.ProgramOccurrence) (*models.ProgramOccurrence, error) {
	result, err := repo.db(ctx).PatchProgramOccurrence(ctx, sqlc.PatchProgramOccurrenceParams{
		StartsAt:    occurrence.StartsAt,
		EndsAt:      occurrence.EndsAt,
		MaxCapacity: occurrence.MaxCapacity,
		ID:          id,
	})
	return database.One(result, err, repo.dbe, mapProgramOccurrence)
}
