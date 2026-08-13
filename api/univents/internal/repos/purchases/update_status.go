package purchases

import (
	"context"
	"errors"

	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (repo *Repo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PurchaseStatus, reason *string) (*models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.UpdateStatus")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdatePurchaseStatus(ctx, sqlc.UpdatePurchaseStatusParams{
		ID:           id,
		Status:       string(status),
		StatusReason: reason,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPurchase(result)), nil
}

// UpdateStatusIf performs a guarded status transition — the row is only
// flipped when it is currently in `from` — and returns (nil, nil) when the
// guard misses. That is the webhook receiver's idempotency mechanism (split
// 4): a duplicate delivery finds the purchase already flipped and becomes a
// no-op instead of re-flipping materialized rows.
func (repo *Repo) UpdateStatusIf(ctx context.Context, id uuid.UUID, from, to models.PurchaseStatus, reason *string) (*models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.UpdateStatusIf")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdatePurchaseStatusIf(ctx, sqlc.UpdatePurchaseStatusIfParams{
		ID:           id,
		FromStatus:   string(from),
		ToStatus:     string(to),
		StatusReason: reason,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//nolint:nilnil // guard missed — the caller treats this as an idempotent no-op
			return nil, nil
		}
		return nil, repo.dbe(err)
	}
	return new(mapPurchase(result)), nil
}

// UpdateRiverJobID links the expiry river job to the purchase (checkout,
// split 7): the job is enqueued with river.InsertTx inside the checkout tx
// and this write stores its id so the webhook receiver can cancel it on
// approve.
func (repo *Repo) UpdateRiverJobID(ctx context.Context, id uuid.UUID, riverJobID int64) (*models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.UpdateRiverJobID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdatePurchaseRiverJob(ctx, sqlc.UpdatePurchaseRiverJobParams{
		ID:         id,
		RiverJobID: &riverJobID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPurchase(result)), nil
}

// AttachIntent stores the Payssage intent on the purchase after the intent
// was created (checkout, split 7, post-commit): seller, intent id (the D2
// correlation key), and the pix QR.
func (repo *Repo) AttachIntent(ctx context.Context, id uuid.UUID, sellerID, intentID uuid.UUID, qrCode, qrCodeBase64 *string) (*models.Purchase, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.AttachIntent")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).AttachIntentToPurchase(ctx, sqlc.AttachIntentToPurchaseParams{
		ID:               id,
		PayssageSellerID: &sellerID,
		PayssageIntentID: &intentID,
		QrCode:           qrCode,
		QrCodeBase64:     qrCodeBase64,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPurchase(result)), nil
}
