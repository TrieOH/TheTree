package infrastructure

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/shared/errx"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type purchaseRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ domain.PurchaseRepository = (*purchaseRepo)(nil)

func NewPurchaseRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) domain.PurchaseRepository {
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

func mapPurchaseFromDB(src *sqlc.Purchase) *domain.Purchase {
	return &domain.Purchase{
		ID:              src.ID,
		EditionID:       src.EditionID,
		UserID:          src.UserID,
		Status:          domain.PurchaseStatus(src.Status),
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

func mapLineItemFromDB(src *sqlc.PurchaseItem) *domain.LineItem {
	return &domain.LineItem{
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

func (repo *purchaseRepo) Create(ctx context.Context, toCreate domain.Purchase) (*domain.Purchase, error) {
	ctx, span := repo.tracer.Start(ctx, "PurchaseRepo.Create")
	defer span.End()

	sqlcPurchase, err := repo.queries(ctx).CreatePurchase(ctx, sqlc.CreatePurchaseParams{
		EditionID:       toCreate.EditionID,
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

func (repo *purchaseRepo) CreateLineItem(ctx context.Context, toCreate domain.LineItem) (*domain.LineItem, error) {
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

func mapTicketGrantFromDB(src *sqlc.GetTicketGrantsByPaymentIntentRow) *domain.TicketGrant {
	return &domain.TicketGrant{
		TicketID: src.TicketID,
		UserID:   src.UserID,
	}
}

func (repo *purchaseRepo) GetTicketIDsByPaymentIntent(ctx context.Context, paymentID string) ([]domain.TicketGrant, error) {
	ctx, span := repo.tracer.Start(ctx, "TicketsRepo.GetTicketIDsByPaymentIntent")
	defer span.End()

	sqlcTicketGrants, err := repo.queries(ctx).GetTicketGrantsByPaymentIntent(ctx, &paymentID)
	if err != nil {
		return nil, errx.FromDB(err, "ticket grant")
	}

	out := make([]domain.TicketGrant, 0, len(sqlcTicketGrants))
	for _, grant := range sqlcTicketGrants {
		out = append(out, *mapTicketGrantFromDB(&grant))
	}
	return out, nil
}
