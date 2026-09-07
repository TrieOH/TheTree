package badges

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EmitForConfirmedRegistration emits the participant badge for a confirmed
// registration and emails it exactly once. Called by the checkout feature when
// a registration is confirmed (paid: webhook approved; free: confirmed at
// inscription). Non-confirmed registrations are a no-op.
func (o *Operations) EmitForConfirmedRegistration(ctx context.Context, registrationID uuid.UUID) (*models.BadgeEmission, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.EmitForConfirmedRegistration")
	defer span.End()

	reg, err := o.registrations.GetByID(ctx, registrationID)
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
			zap.String("registration_id", registrationID.String()))
		//nolint:nilnil // deferred emission is a normal state, not an error
		return nil, nil
	}

	emission, err := o.emissions.Upsert(ctx, &models.BadgeEmission{
		EditionID:      reg.EditionID,
		UserID:         *reg.AttendeeUserID,
		Origin:         models.BadgeEmissionOriginParticipant,
		RegistrationID: &reg.ID,
	})
	if err != nil {
		return nil, err
	}

	if emission.EmailSentAt == nil {
		//nolint:gosec // Detach cancellation AND the caller's transaction before
		// going async: WithoutCancel alone preserves TxKeyValue, so the email
		// goroutine's DB reads would share the caller's tx connection and race
		// it (pgx "conn busy" / SQLSTATE 08P01). The email only needs pool reads.
		go o.sendBadgeEmail(database.WithoutTx(context.WithoutCancel(ctx)), reg, emission)
	}
	return emission, nil
}
