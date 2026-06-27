package purchases

import (
	"context"

	"lib/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueryService struct {
	products  ports.ProductsRepository
	purchases ports.PurchaseRepository
	editions  ports.EditionsRepository
	logger    *zap.Logger
	tracer    trace.Tracer
	tx        database.TxRunner
}

func NewQueryService(
	products ports.ProductsRepository,
	purchases ports.PurchaseRepository,
	editions ports.EditionsRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		products:  products,
		purchases: purchases,
		editions:  editions,
		logger:    logger,
		tracer:    tracer,
		tx:        tx,
	}
}

func (uc *QueryService) GetPurchaseByPaymentID(ctx context.Context, paymentID string) (out *contracts.Purchase, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "PurchaseService.GetPurchaseByPaymentID")
	defer span.End()

	return uc.purchases.GetByPaymentID(ctx, paymentID)
}

func (uc *QueryService) ListUserPurchases(ctx context.Context) (out []contracts.Purchase, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "PurchaseService.ListUserPurchases")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	return uc.purchases.ListUserPurchases(ctx, sub.ID)
}

func (uc *QueryService) ListPurchaseItems(ctx context.Context, purchaseID uuid.UUID) (out []contracts.LineItem, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "PurchaseService.ListPurchaseItems")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	return uc.purchases.ListPurchaseItems(ctx, purchaseID, sub.ID)
}
