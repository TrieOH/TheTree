package badges

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
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

	emission, err := o.emissions.Upsert(ctx, &models.BadgeEmission{
		EditionID:      reg.EditionID,
		UserID:         reg.AttendeeUserID,
		Origin:         models.BadgeEmissionOriginParticipant,
		RegistrationID: &reg.ID,
	})
	if err != nil {
		return nil, err
	}

	if emission.EmailSentAt == nil {
		//nolint:gosec // context.WithoutCancel detaches cancellation while preserving values
		go o.sendBadgeEmail(context.WithoutCancel(ctx), reg, emission)
	}
	return emission, nil
}
