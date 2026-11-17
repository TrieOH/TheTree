package service

import (
	"context"
	"fmt"
	"lib/telemetry"
	"log/slog"
	"time"
	"univents/internal/features/checkouts/jobs"
	"univents/models"
	"uuid"
)

// giftEmailDelay holds the gift email back long enough for the confirming
// tx to commit: it is enqueued inside the tx, before commit.
const giftEmailDelay = 30 * time.Second

// enqueueGiftEmail schedules the gifted-ticket email for an accountless
// free-order recipient inside the checkout tx (the registration is already
// confirmed at checkout). The delay keeps the job from running before the
// confirmation commits; on rollback it finds no registration and no-ops.
// Account holders are skipped: they see the ticket in my-ticket and get
// the badge email.
func (o *Operations) enqueueGiftEmail(ctx context.Context, reg *models.Registration) {
	if reg.AttendeeUserID != nil {
		return
	}
	err := o.Jobs.EnqueueAfter(ctx, jobs.SendGiftEmailArgs{RegistrationID: reg.ID}, giftEmailDelay)
	if err != nil {
		telemetry.Log().Error("checkout: failed to enqueue gift email",
			slog.String("registration_id", reg.ID.String()),
			telemetry.Err(err))
	}
}

// enqueueExpiry enqueues the purchases.expire job inside the checkout tx,
// before commit, delayed to the purchase's expires_at. The delay means it
// cannot run before the purchase is visible; an enqueue failure fails the
// checkout; a rolled-back checkout leaves a job that finds no purchase and
// no-ops.
func (o *Operations) enqueueExpiry(ctx context.Context, purchaseID uuid.UUID, expiresAt time.Time) error {
	err := o.Jobs.EnqueueAfter(ctx, jobs.ExpirePurchaseArgs{PurchaseID: purchaseID}, time.Until(expiresAt))
	if err != nil {
		return fmt.Errorf("checkout: enqueue expiry: %w", err)
	}
	return nil
}
