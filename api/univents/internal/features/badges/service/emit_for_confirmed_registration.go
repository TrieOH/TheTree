package service

import (
	"context"
	"lib/telemetry"
	"log/slog"
	"time"
	"univents/internal/features/badges/jobs"
	"univents/models"
	"uuid"
)

// badgeEmailDelay holds the badge email back long enough for the caller's tx
// to commit: it is enqueued inside the tx, before commit. On rollback the job
// finds no confirmed registration and no-ops.
const badgeEmailDelay = 30 * time.Second

// EmitForConfirmedRegistration emits the participant badge for a confirmed
// registration and enqueues its email (the badges.send_email job sends it
// exactly once). Called by the checkout feature when a registration is
// confirmed (paid: webhook approved; free: confirmed at inscription).
// Non-confirmed registrations are a no-op.
func (o *Operations) EmitForConfirmedRegistration(ctx context.Context, registrationID uuid.UUID) (*models.BadgeEmission, error) {
	reg, err := o.Registrations.GetByID(ctx, registrationID)
	if err != nil {
		return nil, err
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		//nolint:nilnil // no-op contract: pending/failed registrations simply never emit
		return nil, nil
	}
	if reg.AttendeeUserID == nil {
		// Accountless gifted ticket (email-only recipient): there is no
		// profile to attach the badge to or email it to (the badge email's
		// QR points at the profile). Emission is deferred until the
		// recipient claims an account and the registration gets its actor
		// id — the claim flow re-emits from here.
		telemetry.Log().Info("badge emission deferred: attendee has no account yet",
			slog.String("registration_id", registrationID.String()))
		//nolint:nilnil // deferred emission is a normal state, not an error
		return nil, nil
	}

	emission, err := o.Emissions.Upsert(ctx, &models.BadgeEmission{
		EditionID:      reg.EditionID,
		UserID:         *reg.AttendeeUserID,
		Origin:         models.BadgeEmissionOriginParticipant,
		RegistrationID: &reg.ID,
	})
	if err != nil {
		return nil, err
	}

	if emission.EmailSentAt == nil {
		err = o.Jobs.EnqueueAfter(ctx, jobs.SendBadgeEmailArgs{EmissionID: emission.ID, RegistrationID: reg.ID}, badgeEmailDelay)
		if err != nil {
			telemetry.Log().Error("failed to enqueue badge email",
				slog.String("registration_id", reg.ID.String()),
				telemetry.Err(err))
		}
	}
	return emission, nil
}
