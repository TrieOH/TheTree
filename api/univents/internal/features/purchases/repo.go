package purchases

import (
	"context"

	"univents/internal/platform/database"
	"univents/internal/platform/database/sqlc"
	"univents/internal/shared/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type purchaseRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.PurchaseRepository = (*purchaseRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.PurchaseRepository {
	return &purchaseRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *purchaseRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapPurchaseFromDB(src *sqlc.Purchase) *contracts.Purchase {
	return &contracts.Purchase{
		ID:              src.ID,
		EditionID:       src.EditionID,
		SessionID:       src.SessionID,
		UserID:          src.UserID,
		Status:          contracts.PurchaseStatus(src.Status),
		SubtotalCents:   src.SubtotalCents,
		DiscountCents:   src.DiscountCents,
		TaxCents:        src.TaxCents,
		TaxBreakdown:    src.TaxBreakdown,
		TotalCents:      src.TotalCents,
		PaymentProvider: src.PaymentProvider,
		PaymentID:       src.PaymentID,
		FulfilledAt:     src.FulfilledAt,
		FulfilmentNotes: src.FulfilmentNotes,
		CreatedAt:       src.CreatedAt,
		UpdatedAt:       src.UpdatedAt,
		DeletedAt:       src.DeletedAt,
	}
}

func mapLineItemFromDB(src *sqlc.PurchaseItem) *contracts.LineItem {
	return &contracts.LineItem{
		ID:                  src.ID,
		PurchaseID:          src.PurchaseID,
		ItemType:            src.ItemType,
		ItemID:              src.ItemID,
		Quantity:            src.Quantity,
		UnitPriceCents:      src.UnitPriceCents,
		TotalPriceCents:     src.TotalPriceCents,
		AssignedToUserID:    src.AssignedToUserID,
		Fulfilled:           src.Fulfilled,
		FulfilledAt:         src.FulfilledAt,
		RefundedQuantity:    src.RefundedQuantity,
		RefundedAmountCents: src.RefundedAmountCents,
		CreatedAt:           src.CreatedAt,
	}
}

func (repo *purchaseRepo) Create(ctx context.Context, toCreate contracts.Purchase) (*contracts.Purchase, error) {
	ctx, span := repo.tracer.Start(ctx, "PurchaseRepo.Create")
	defer span.End()

	sqlcPurchase, err := repo.queries(ctx).CreatePurchase(ctx, sqlc.CreatePurchaseParams{
		EditionID:       toCreate.EditionID,
		SessionID:       toCreate.SessionID,
		UserID:          toCreate.UserID,
		Status:          sqlc.PurchaseStatus(toCreate.Status),
		SubtotalCents:   toCreate.SubtotalCents,
		TotalCents:      toCreate.TotalCents,
		PaymentProvider: toCreate.PaymentProvider,
		PaymentID:       toCreate.PaymentID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "purchase")
	}

	return mapPurchaseFromDB(&sqlcPurchase), nil
}

func (repo *purchaseRepo) GetByPaymentID(ctx context.Context, paymentID string) (*contracts.Purchase, error) {
	ctx, span := repo.tracer.Start(ctx, "PurchaseRepo.GetByPaymentID")
	defer span.End()

	sqlcPurchase, err := repo.queries(ctx).GetPurchaseByPaymentID(ctx, &paymentID)
	if err != nil {
		return nil, errx.FromDB(err, "purchase")
	}

	return mapPurchaseFromDB(&sqlcPurchase), nil
}

func (repo *purchaseRepo) GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*contracts.Purchase, error) {
	ctx, span := repo.tracer.Start(ctx, "PurchaseRepo.GetBySessionID")
	defer span.End()

	sqlcPurchase, err := repo.queries(ctx).GetPurchaseBySessionID(ctx, &sessionID)
	if err != nil {
		return nil, errx.FromDB(err, "purchase")
	}

	return mapPurchaseFromDB(&sqlcPurchase), nil
}

func (repo *purchaseRepo) CreateLineItem(ctx context.Context, toCreate contracts.LineItem) (*contracts.LineItem, error) {
	ctx, span := repo.tracer.Start(ctx, "PurchaseRepo.CreateLineItem")
	defer span.End()

	sqlcLineItem, err := repo.queries(ctx).CreatePurchaseItem(ctx, sqlc.CreatePurchaseItemParams{
		PurchaseID:      toCreate.PurchaseID,
		ItemType:        toCreate.ItemType,
		ItemID:          toCreate.ItemID,
		Quantity:        toCreate.Quantity,
		UnitPriceCents:  toCreate.UnitPriceCents,
		TotalPriceCents: toCreate.TotalPriceCents,
	})
	if err != nil {
		return nil, errx.FromDB(err, "purchase")
	}

	return mapLineItemFromDB(&sqlcLineItem), nil
}

func (repo *purchaseRepo) ConfirmPurchase(ctx context.Context, paymentID string) error {
	ctx, span := repo.tracer.Start(ctx, "PurchaseRepo.ConfirmPurchase")
	defer span.End()

	err := repo.queries(ctx).ConfirmPurchase(ctx, &paymentID)
	if err != nil {
		return errx.FromDB(err, "purchase")
	}

	return nil
}

func (repo *purchaseRepo) CancelPurchase(ctx context.Context, paymentID string) error {
	ctx, span := repo.tracer.Start(ctx, "PurchaseRepo.CancelPurchase")
	defer span.End()

	err := repo.queries(ctx).CancelPurchase(ctx, &paymentID)
	if err != nil {
		return errx.FromDB(err, "purchase")
	}

	return nil
}

func mapTicketGrantFromDB(src *sqlc.GetTicketGrantsByPaymentIntentRow) *contracts.TicketGrant {
	return &contracts.TicketGrant{
		TicketID: src.TicketID,
		UserID:   src.UserID,
	}
}

func (repo *purchaseRepo) GetTicketIDsByPaymentIntent(ctx context.Context, paymentID string) ([]contracts.TicketGrant, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketsRepo.GetTicketIDsByPaymentIntent")
	defer span.End()

	sqlcTicketGrants, err := repo.queries(ctx).GetTicketGrantsByPaymentIntent(ctx, &paymentID)
	if err != nil {
		return nil, errx.FromDB(err, "ticket grant")
	}

	out := make([]contracts.TicketGrant, 0, len(sqlcTicketGrants))
	for _, grant := range sqlcTicketGrants {
		out = append(out, *mapTicketGrantFromDB(&grant))
	}
	return out, nil
}

func (repo *purchaseRepo) ListUserPurchases(ctx context.Context, userID uuid.UUID) ([]contracts.Purchase, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.ListUserPurchases")
	defer span.End()

	sqlcPurchases, err := repo.queries(ctx).ListUserPurchases(ctx, userID)
	if err != nil {
		return nil, errx.FromDB(err, "purchase")
	}

	out := make([]contracts.Purchase, 0, len(sqlcPurchases))
	for _, purchase := range sqlcPurchases {
		out = append(out, *mapPurchaseFromDB(&purchase))
	}
	return out, nil
}

func (repo *purchaseRepo) ListPurchaseItems(ctx context.Context, purchaseID, userID uuid.UUID) ([]contracts.LineItem, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.ListPurchaseItems")
	defer span.End()

	sqlcItems, err := repo.queries(ctx).ListPurchaseItems(ctx, sqlc.ListPurchaseItemsParams{
		PurchaseID: purchaseID,
		UserID:     userID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "purchase item")
	}

	out := make([]contracts.LineItem, 0, len(sqlcItems))
	for _, item := range sqlcItems {
		out = append(out, *mapLineItemFromDB(&item))
	}
	return out, nil
}
