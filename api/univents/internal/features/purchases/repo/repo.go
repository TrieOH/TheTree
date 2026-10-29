package repo

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

var _ ports.PurchaseRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrPurchaseNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapPurchase(src sqlc.Purchase) models.Purchase {
	return models.Purchase{
		ID:               src.ID,
		EditionID:        src.EditionID,
		PurchaserID:      src.PurchaserID,
		Status:           models.PurchaseStatus(src.Status),
		StatusReason:     src.StatusReason,
		TotalCents:       src.TotalCents,
		Currency:         src.Currency,
		PaymentMethod:    src.PaymentMethod,
		PayssageSellerID: src.PayssageSellerID,
		PayssageIntentID: src.PayssageIntentID,
		PayerEmail:       src.PayerEmail,
		QRCode:           src.QrCode,
		QRCodeBase64:     src.QrCodeBase64,
		ExpiresAt:        src.ExpiresAt,
		RiverJobID:       src.RiverJobID,
		CreatedAt:        src.CreatedAt,
		UpdatedAt:        src.UpdatedAt,
		DeletedAt:        src.DeletedAt,
	}
}

func mapPurchaseItem(src sqlc.PurchaseItem) models.PurchaseItem {
	return models.PurchaseItem{
		ID:                src.ID,
		PurchaseID:        src.PurchaseID,
		ItemType:          models.PurchaseItemType(src.ItemType),
		ItemID:            src.ItemID,
		Quantity:          src.Quantity,
		UnitPriceCents:    src.UnitPriceCents,
		RegistrationID:    src.RegistrationID,
		ProductPurchaseID: src.ProductPurchaseID,
		ParticipationID:   src.ParticipationID,
	}
}

// Availability returns the stock position of every purchasable item in the
// edition (available = base - reserved; nil base = unlimited). It fans out
// to the three per-type queries so callers can run it inside the checkout
// tx under the item row locks (split 7); item ids never cross tables.
func (repo *Repo) Availability(ctx context.Context, editionID uuid.UUID) ([]models.ItemAvailability, error) {
	var out []models.ItemAvailability

	tickets, err := repo.db(ctx).ListTicketTypeAvailability(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	for _, row := range tickets {
		out = append(out, models.ItemAvailability{
			ItemType:         models.PurchaseItemTypeTicket,
			ItemID:           row.ItemID,
			BaseQuantity:     row.BaseQuantity,
			ReservedQuantity: row.ReservedQuantity,
		})
	}

	variants, err := repo.db(ctx).ListProductVariantAvailability(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	for _, row := range variants {
		out = append(out, models.ItemAvailability{
			ItemType:         models.PurchaseItemTypeProduct,
			ItemID:           row.ItemID,
			BaseQuantity:     row.BaseQuantity,
			ReservedQuantity: row.ReservedQuantity,
		})
	}

	occurrences, err := repo.db(ctx).ListProgramOccurrenceAvailability(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	for _, row := range occurrences {
		out = append(out, models.ItemAvailability{
			ItemType:         models.PurchaseItemTypeProgramOccurrence,
			ItemID:           row.ItemID,
			BaseQuantity:     row.BaseQuantity,
			ReservedQuantity: row.ReservedQuantity,
		})
	}

	return out, nil
}

func (repo *Repo) CreatePurchase(ctx context.Context, toCreate *models.Purchase) (*models.Purchase, error) {
	result, err := repo.db(ctx).CreatePurchase(ctx, sqlc.CreatePurchaseParams{
		EditionID:        toCreate.EditionID,
		PurchaserID:      toCreate.PurchaserID,
		Status:           string(toCreate.Status),
		StatusReason:     toCreate.StatusReason,
		TotalCents:       toCreate.TotalCents,
		Currency:         toCreate.Currency,
		PaymentMethod:    toCreate.PaymentMethod,
		PayssageSellerID: toCreate.PayssageSellerID,
		PayssageIntentID: toCreate.PayssageIntentID,
		PayerEmail:       toCreate.PayerEmail,
		QrCode:           toCreate.QRCode,
		QrCodeBase64:     toCreate.QRCodeBase64,
		ExpiresAt:        toCreate.ExpiresAt,
		RiverJobID:       toCreate.RiverJobID,
	})
	return database.One(result, err, repo.dbe, mapPurchase)
}

func (repo *Repo) CreatePurchaseItem(ctx context.Context, toCreate *models.PurchaseItem) (*models.PurchaseItem, error) {
	result, err := repo.db(ctx).CreatePurchaseItem(ctx, sqlc.CreatePurchaseItemParams{
		PurchaseID:        toCreate.PurchaseID,
		ItemType:          string(toCreate.ItemType),
		ItemID:            toCreate.ItemID,
		Quantity:          toCreate.Quantity,
		UnitPriceCents:    toCreate.UnitPriceCents,
		RegistrationID:    toCreate.RegistrationID,
		ProductPurchaseID: toCreate.ProductPurchaseID,
		ParticipationID:   toCreate.ParticipationID,
	})
	return database.One(result, err, repo.dbe, mapPurchaseItem)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Purchase, error) {
	result, err := repo.db(ctx).GetPurchaseByID(ctx, id)
	return database.One(result, err, repo.dbe, mapPurchase)
}

func (repo *Repo) GetByIDForOwner(ctx context.Context, id, purchaserID uuid.UUID) (*models.Purchase, error) {
	result, err := repo.db(ctx).GetPurchaseByIDForOwner(ctx, sqlc.GetPurchaseByIDForOwnerParams{
		ID:          id,
		PurchaserID: purchaserID,
	})
	return database.One(result, err, repo.dbe, mapPurchase)
}

func (repo *Repo) GetByIntentID(ctx context.Context, intentID uuid.UUID) (*models.Purchase, error) {
	result, err := repo.db(ctx).GetPurchaseByIntentID(ctx, &intentID)
	return database.One(result, err, repo.dbe, mapPurchase)
}

func (repo *Repo) ListByPurchaser(ctx context.Context, purchaserID uuid.UUID) ([]models.Purchase, error) {
	rows, err := repo.db(ctx).ListPurchasesByPurchaser(ctx, purchaserID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.Purchase, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapPurchase(row))
	}
	return out, nil
}

func (repo *Repo) ListItemsByPurchase(ctx context.Context, purchaseID uuid.UUID) ([]models.PurchaseItem, error) {
	rows, err := repo.db(ctx).ListPurchaseItemsByPurchase(ctx, purchaseID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.PurchaseItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapPurchaseItem(row))
	}
	return out, nil
}

// ListByEdition returns every purchase of an edition, newest first — the
// organizer orders read (refund plan B3).
func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Purchase, error) {
	rows, err := repo.db(ctx).ListPurchasesByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.Purchase, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapPurchase(row))
	}
	return out, nil
}

func (repo *Repo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PurchaseStatus, reason *string) (*models.Purchase, error) {
	result, err := repo.db(ctx).UpdatePurchaseStatus(ctx, sqlc.UpdatePurchaseStatusParams{
		ID:           id,
		Status:       string(status),
		StatusReason: reason,
	})
	return database.One(result, err, repo.dbe, mapPurchase)
}

// UpdateStatusIf performs a guarded status transition — the row is only
// flipped when it is currently in `from` — and returns (nil, nil) when the
// guard misses. That is the webhook receiver's idempotency mechanism (split
// 4): a duplicate delivery finds the purchase already flipped and becomes a
// no-op instead of re-flipping materialized rows.
func (repo *Repo) UpdateStatusIf(ctx context.Context, id uuid.UUID, from, to models.PurchaseStatus, reason *string) (*models.Purchase, error) {
	result, err := repo.db(ctx).UpdatePurchaseStatusIf(ctx, sqlc.UpdatePurchaseStatusIfParams{
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
	result, err := repo.db(ctx).UpdatePurchaseRiverJob(ctx, sqlc.UpdatePurchaseRiverJobParams{
		ID:         id,
		RiverJobID: &riverJobID,
	})
	return database.One(result, err, repo.dbe, mapPurchase)
}

// AttachIntent stores the Payssage intent on the purchase after the intent
// was created (checkout, split 7, post-commit): seller, intent id (the D2
// correlation key), and the pix QR.
func (repo *Repo) AttachIntent(ctx context.Context, id uuid.UUID, sellerID, intentID uuid.UUID, qrCode, qrCodeBase64 *string) (*models.Purchase, error) {
	result, err := repo.db(ctx).AttachIntentToPurchase(ctx, sqlc.AttachIntentToPurchaseParams{
		ID:               id,
		PayssageSellerID: &sellerID,
		PayssageIntentID: &intentID,
		QrCode:           qrCode,
		QrCodeBase64:     qrCodeBase64,
	})
	return database.One(result, err, repo.dbe, mapPurchase)
}
