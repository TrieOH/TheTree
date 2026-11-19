// Package programs is the programs feature: its repo, handlers and
// operations. Build wires it.
package programs

import (
	"context"
	"encoding/json"
	"errors"
	"lib/database"
	"lib/telemetry"
	"log/slog"
	"univents/errs"
	"univents/internal/authz"
	"univents/internal/notify"
	"univents/models"
	"univents/ports"
	"uuid"
)

// Notifier is the LISTEN/NOTIFY publish surface for the store's stock
// deltas (register/de-register move occurrence availability). Satisfied by
// *database.Notifier (lib/go/database).
type Notifier interface {
	Notify(ctx context.Context, channel, payload string) error
}

// channelUniventsChanges is the single NOTIFY channel the store publishes
// on (contract in internal/notify).
const channelUniventsChanges = notify.Channel

type Deps struct {
	Events         ports.EventRepo
	Editions       ports.EditionRepo
	Programs       ports.ProgramRepo
	Occurrences    ports.ProgramOccurrenceRepo
	Registrations  ports.RegistrationRepo
	TicketTypes    ports.TicketTypeRepo
	Participations ports.ProgramParticipationRepo
	Authz          *authz.Service
	Notifier       Notifier
	Tx             database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// notifyStock publishes the occurrence's stock delta (D10 contract): item
// ids only — the SSE relay re-queries availability from the DB. A missed
// notification is a stale snapshot, never data loss.
func (o *Operations) notifyStock(ctx context.Context, editionID, occurrenceID uuid.UUID) {
	raw, err := json.Marshal(notify.StockNotification{
		Kind:      notify.KindStock,
		EditionID: editionID,
		ItemIDs:   []uuid.UUID{occurrenceID},
	})
	if err != nil {
		telemetry.Log().Error("programs: marshal notifier payload",
			slog.String("channel", channelUniventsChanges),
			telemetry.Err(err))
		return
	}
	err = o.Notifier.Notify(ctx, channelUniventsChanges, string(raw))
	if err != nil {
		telemetry.Log().Error("programs: publish to notifier",
			slog.String("channel", channelUniventsChanges),
			slog.String("payload", string(raw)),
			telemetry.Err(err))
	}
}

// CheckIn marks an attendee as present at a checkpoint occurrence
// (staff-only). Checkpoints are pass-through by design — attendees never
// register for them — so the gate is "can they pass", not "did they sign
// up": a confirmed edition ticket plus the checkpoint's access rules
// (staff_only, min_access_level), resolved from the attendee's badge QR.
// The first scan upserts an attended participation anchored to the
// attendee's edition-ticket registration; re-scanning is idempotent (an
// existing live row is flipped to attended, never duplicated — the
// partial unique index uniq_program_participations_active_per_occurrence_attendee
// is the backstop). Activities keep the participation-driven
// mark-attended flow and are rejected here (400).
func (o *Operations) CheckIn(ctx context.Context, occurrenceID, attendeeID, actorID uuid.UUID) (*models.ProgramParticipation, error) {
	occ, err := o.Occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, occ.EditionID)
	if err != nil {
		return nil, err
	}
	err = o.Authz.EventFor(ctx, actorID, edition.EventID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	program, err := o.Programs.GetByID(ctx, occ.ProgramID)
	if err != nil {
		return nil, err
	}
	if program.Kind != models.ProgramKindCheckpoint {
		return nil, errs.ErrCheckInNotCheckpoint
	}

	reg, err := o.Registrations.GetActiveByEditionAndAttendee(ctx, occ.EditionID, attendeeID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errs.ErrTicketRequiredForCheckpoint
		}
		return nil, err
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		return nil, errs.ErrTicketRequiredForCheckpoint
	}
	err = o.checkProgramAccess(ctx, occ, program, reg, attendeeID)
	if err != nil {
		return nil, err
	}

	return o.Participations.UpsertAttended(ctx, occ.EditionID, occurrenceID, reg.ID)
}

func (o *Operations) CreateOccurrence(ctx context.Context, payload models.CreateProgramOccurrenceInput) (*models.ProgramOccurrence, error) {
	program, err := o.Programs.GetByID(ctx, payload.ProgramID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, program.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	occurrence := &models.ProgramOccurrence{
		ProgramID:   payload.ProgramID,
		EditionID:   program.EditionID,
		StartsAt:    payload.StartsAt,
		EndsAt:      payload.EndsAt,
		MaxCapacity: payload.MaxCapacity,
	}

	return o.Occurrences.CreateOccurrence(ctx, occurrence)
}

func (o *Operations) CreateProgram(ctx context.Context, payload models.CreateProgramInput) (*models.Program, error) {
	edition, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	program := &models.Program{
		EditionID:      payload.EditionID,
		Kind:           payload.Kind,
		Name:           payload.Name,
		Description:    payload.Description,
		MinAccessLevel: payload.MinAccessLevel,
		StaffOnly:      payload.StaffOnly,
		Price:          payload.Price,
	}

	return o.Programs.Create(ctx, program)
}

func (o *Operations) DeleteOccurrence(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	existing, err := o.Occurrences.GetOccurrenceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Occurrences.DeleteOccurrence(ctx, id)
}

func (o *Operations) DeleteProgram(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	existing, err := o.Programs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Programs.Delete(ctx, id)
}

// Deregister cancels the caller's sign-up for an occurrence (self-service).
// The spot is locked once attendance is marked: the guarded update flips
// only 'registered' rows, so an attended/no_show participation (or a racing
// state change) returns 0 rows → 409. Deregistering is a status flip, never
// a delete — the participation stays in the ledger as cancelled history.
func (o *Operations) Deregister(ctx context.Context, occurrenceID, userID uuid.UUID) (*models.ProgramParticipation, error) {
	occ, err := o.Occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	reg, err := o.Registrations.GetActiveByEditionAndAttendee(ctx, occ.EditionID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errs.ErrNotRegistered
		}
		return nil, err
	}

	part, err := o.Participations.GetActiveByOccurrenceAndRegistration(ctx, occurrenceID, reg.ID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errs.ErrNotRegistered
		}
		return nil, err
	}

	updated, err := o.Participations.UpdateParticipationStatusIfRegistered(ctx, part.ID, models.ProgramParticipationStatusCancelled)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errs.ErrRegistrationLocked
	}

	o.notifyStock(ctx, occ.EditionID, occurrenceID)
	return updated, nil
}

func (o *Operations) GetOccurrenceByID(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	return o.Occurrences.GetOccurrenceByID(ctx, id)
}

func (o *Operations) GetProgramByID(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	return o.Programs.GetByID(ctx, id)
}

func (o *Operations) ListOccurrencesByProgram(ctx context.Context, programID uuid.UUID) ([]models.ProgramOccurrence, error) {
	return o.Occurrences.ListOccurrencesByProgram(ctx, programID)
}

func (o *Operations) ListOccurrencesByEdition(ctx context.Context, editionID uuid.UUID) ([]models.ProgramOccurrence, error) {
	return o.Occurrences.ListOccurrencesByEdition(ctx, editionID)
}

func (o *Operations) ListProgramsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Program, error) {
	return o.Programs.ListByEdition(ctx, editionID)
}

// MarkAttended flips a participation to attended (staff-only). Any event
// staff (owner/admin/staff) may mark any activity in the edition — the
// domain's per-program staff assignment is deferred. The guarded update
// flips any non-cancelled row: re-marking an attended participation is
// idempotent, a no_show→attended correction is allowed, and a cancelled
// participation (refunded purchase, or the attendee dropped out) is 409.
func (o *Operations) MarkAttended(ctx context.Context, participationID, actorID uuid.UUID) (*models.ProgramParticipation, error) {
	part, err := o.Participations.GetParticipationByID(ctx, participationID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, part.EditionID)
	if err != nil {
		return nil, err
	}
	err = o.Authz.EventFor(ctx, actorID, edition.EventID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	updated, err := o.Participations.UpdateParticipationStatusIfNotCancelled(ctx, participationID, models.ProgramParticipationStatusAttended)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errs.ErrParticipationCancelled
	}
	return updated, nil
}

// MyParticipations returns the caller's live activity sign-ups in the
// edition ("my activities"), joined with program + occurrence. Returns
// (nil, nil) when the caller holds no ticket — no sign-ups is a normal
// state, like the my-ticket read. Cancelled rows are history and never
// shown: the append-only ledger's current-state projection.
func (o *Operations) MyParticipations(ctx context.Context, editionID, userID uuid.UUID) ([]models.MyParticipation, error) {
	_, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	reg, err := o.Registrations.GetActiveByEditionAndAttendee(ctx, editionID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// no ticket → no sign-ups is a normal state, not an error
			return nil, nil
		}
		return nil, err
	}

	return o.Participations.ListActiveByEditionAndRegistration(ctx, editionID, reg.ID)
}

// Participants returns the attendance surface of an occurrence: every live
// participation with the attendee's identity (staff-only, owner/admin/staff
// — the same role the mark-attended write requires).
func (o *Operations) Participants(ctx context.Context, occurrenceID, actorID uuid.UUID) ([]models.ProgramParticipant, error) {
	occ, err := o.Occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, occ.EditionID)
	if err != nil {
		return nil, err
	}
	err = o.Authz.EventFor(ctx, actorID, edition.EventID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	return o.Participations.ListByOccurrence(ctx, occurrenceID)
}

func (o *Operations) PatchOccurrence(ctx context.Context, payload models.PatchProgramOccurrenceInput) (*models.ProgramOccurrence, error) {
	existing, err := o.Occurrences.GetOccurrenceByID(ctx, payload.OccurrenceID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	occurrence := &models.ProgramOccurrence{
		StartsAt:    payload.StartsAt,
		EndsAt:      payload.EndsAt,
		MaxCapacity: payload.MaxCapacity,
	}

	return o.Occurrences.PatchOccurrence(ctx, payload.OccurrenceID, occurrence)
}

func (o *Operations) PatchProgram(ctx context.Context, payload models.PatchProgramInput) (*models.Program, error) {
	existing, err := o.Programs.GetByID(ctx, payload.ProgramID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	program := &models.Program{
		Kind:           payload.Kind,
		Name:           payload.Name,
		Description:    payload.Description,
		MinAccessLevel: payload.MinAccessLevel,
		StaffOnly:      payload.StaffOnly,
		Price:          payload.Price,
		BannerURL:      payload.BannerURL,
	}

	return o.Programs.Patch(ctx, payload.ProgramID, program)
}

// checkEligible resolves the caller's confirmed registration in the
// occurrence's edition and enforces the activity's access rules (D1):
//
//   - a confirmed (paid or free) ticket is required — pending/unpaid
//     reservations cannot hold activity spots;
//   - staff_only programs additionally require an event staff role
//     (owner/admin/staff);
//   - min_access_level gates on the ticket's access level (null = open).
//
// Returns the registration and program, or a typed fun error (403).
func (o *Operations) checkEligible(ctx context.Context, occ *models.ProgramOccurrence, userID uuid.UUID) (*models.Registration, *models.Program, error) {
	reg, err := o.Registrations.GetActiveByEditionAndAttendee(ctx, occ.EditionID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, nil, errs.ErrTicketRequiredForActivity
		}
		return nil, nil, err
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		return nil, nil, errs.ErrTicketRequiredForActivity
	}

	program, err := o.Programs.GetByID(ctx, occ.ProgramID)
	if err != nil {
		return nil, nil, err
	}

	err = o.checkProgramAccess(ctx, occ, program, reg, userID)
	if err != nil {
		return nil, nil, err
	}
	return reg, program, nil
}

// checkProgramAccess enforces the program's access rules for an attendee
// (shared by register and checkpoint check-in):
//
//   - staff_only programs additionally require an event staff role
//     (owner/admin/staff);
//   - min_access_level gates on the ticket's access level (null = open).
func (o *Operations) checkProgramAccess(ctx context.Context, occ *models.ProgramOccurrence, program *models.Program, reg *models.Registration, userID uuid.UUID) error {
	if program.StaffOnly {
		edition, err := o.Editions.GetByID(ctx, occ.EditionID)
		if err != nil {
			return err
		}
		err = o.Authz.EventFor(ctx, userID, edition.EventID, models.EventMemberRoleStaff)
		if err != nil {
			return err
		}
	}

	if program.MinAccessLevel != nil {
		ticket, err := o.TicketTypes.GetByID(ctx, reg.TicketTypeID)
		if err != nil {
			return err
		}
		if ticket.AccessLevel < *program.MinAccessLevel {
			return &errs.AccessLevelError{Kind: string(program.Kind)}
		}
	}

	return nil
}

// Register signs the caller up for an occurrence (self-service). Flow:
// resolve + permission-check (confirmed ticket, staff_only, access level),
// then inside a tx: row-lock the occurrence (serializes concurrent
// sign-ups), 409 on double registration, 409 when full, insert. The partial
// unique index
// (uniq_program_participations_active_per_occurrence_attendee) is the
// atomic backstop for two racing inserts. On success the occurrence's stock
// delta is published so the store's realtime surfaces stay fresh.
func (o *Operations) Register(ctx context.Context, occurrenceID, userID uuid.UUID) (*models.ProgramParticipation, error) {
	occ, err := o.Occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	reg, _, err := o.checkEligible(ctx, occ, userID)
	if err != nil {
		return nil, err
	}

	var part *models.ProgramParticipation
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		// Lock the occurrence row: two concurrent sign-ups serialize here,
		// so the count below cannot oversell the last slot.
		_, err := o.Occurrences.GetOccurrenceByIDForUpdate(ctx, occurrenceID)
		if err != nil {
			return err
		}

		existing, err := o.Participations.GetActiveByOccurrenceAndRegistration(ctx, occurrenceID, reg.ID)
		if existing != nil {
			return errs.ErrAlreadyRegistered
		}
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return err
		}

		count, err := o.Participations.CountActiveByOccurrence(ctx, occurrenceID)
		if err != nil {
			return err
		}
		if occ.MaxCapacity != nil && count >= int64(*occ.MaxCapacity) {
			return errs.ErrOccurrenceFull
		}

		part, err = o.Participations.CreateParticipation(ctx, &models.ProgramParticipation{
			EditionID:      occ.EditionID,
			OccurrenceID:   occurrenceID,
			RegistrationID: reg.ID,
			Status:         models.ProgramParticipationStatusRegistered,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	o.notifyStock(ctx, occ.EditionID, occurrenceID)
	return part, nil
}
