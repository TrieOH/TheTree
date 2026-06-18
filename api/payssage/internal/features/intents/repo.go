package intents

import (
	"context"
	"encoding/json"
	"payssage/ports"

	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/internal/shared/errx"
	"payssage/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type intentsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.IntentRepository = (*intentsRepo)(nil)

func NewIntentsRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.IntentRepository {
	return &intentsRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *intentsRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapIntentFromDB(src *sqlc.Intent) (*models.Intent, error) {
	intent := &models.Intent{
		ID:          src.ID,
		WorkspaceID: src.WorkspaceID,
		Amount:      src.Amount,
		Currency:    src.Currency,
		Status:      models.IntentStatus(src.Status),
		Provider:    src.Provider,
		Metadata:    src.Metadata,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}

	if src.SellerCredentialID != nil {
		intent.SellerCredentialID = src.SellerCredentialID
	}

	switch src.Provider {
	case "mercadopago":
		if src.ProviderData != nil {
			var mp models.MercadoPagoIntentData
			if err := json.Unmarshal(src.ProviderData, &mp); err != nil {
				return nil, errx.Internal("intent").SetMessage("failed to unmarshal mercadopago provider data").SetCause(err)
			}
			intent.MercadoPagoData = &mp
		}
	}

	return intent, nil
}

func marshalProviderData(intent models.Intent) (json.RawMessage, error) {
	switch intent.Provider {
	case "mercadopago":
		if intent.MercadoPagoData == nil {
			return json.RawMessage("{}"), nil
		}
		b, err := json.Marshal(intent.MercadoPagoData)
		if err != nil {
			return nil, errx.Internal("intent").SetMessage("failed to marshal mercadopago provider data").SetCause(err)
		}
		return b, nil
	default:
		return json.RawMessage("{}"), nil
	}
}

func (repo *intentsRepo) Create(ctx context.Context, toCreate models.Intent) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.Create")
	defer span.End()

	providerData, err := marshalProviderData(toCreate)
	if err != nil {
		return nil, err
	}

	sqlcIntent, err := repo.queries(ctx).CreateIntent(ctx, sqlc.CreateIntentParams{
		ID:                 toCreate.ID,
		WorkspaceID:        toCreate.WorkspaceID,
		Amount:             toCreate.Amount,
		Currency:           toCreate.Currency,
		Status:             sqlc.IntentStatus(toCreate.Status),
		Provider:           toCreate.Provider,
		ProviderData:       providerData,
		Metadata:           toCreate.Metadata,
		SellerCredentialID: toCreate.SellerCredentialID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	return mapIntentFromDB(&sqlcIntent)
}

func (repo *intentsRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.GetByID")
	defer span.End()

	sqlcIntent, err := repo.queries(ctx).GetIntentByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	return mapIntentFromDB(&sqlcIntent)
}

func (repo *intentsRepo) List(ctx context.Context) ([]models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.List")
	defer span.End()

	sqlcIntents, err := repo.queries(ctx).ListIntents(ctx)
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	out := make([]models.Intent, 0, len(sqlcIntents))
	for _, intent := range sqlcIntents {
		mapped, err := mapIntentFromDB(&intent)
		if err != nil {
			return nil, err
		}
		out = append(out, *mapped)
	}
	return out, nil
}

func (repo *intentsRepo) ListIntentsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.ListIntentsByWorkspace")
	defer span.End()

	sqlcIntents, err := repo.queries(ctx).ListIntentsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	out := make([]models.Intent, 0, len(sqlcIntents))
	for _, intent := range sqlcIntents {
		mapped, err := mapIntentFromDB(&intent)
		if err != nil {
			return nil, err
		}
		out = append(out, *mapped)
	}
	return out, nil
}

func (repo *intentsRepo) Cancel(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.Cancel")
	defer span.End()

	sqlcIntent, err := repo.queries(ctx).CancelIntent(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	return mapIntentFromDB(&sqlcIntent)
}

func (repo *intentsRepo) Confirm(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.Confirm")
	defer span.End()

	sqlcIntent, err := repo.queries(ctx).ConfirmIntent(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	return mapIntentFromDB(&sqlcIntent)
}

func (repo *intentsRepo) Fail(ctx context.Context, id uuid.UUID) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.Fail")
	defer span.End()

	sqlcIntent, err := repo.queries(ctx).FailIntent(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	return mapIntentFromDB(&sqlcIntent)
}

func (repo *intentsRepo) UpdateProviderData(ctx context.Context, intent models.Intent) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.UpdateProviderData")
	defer span.End()

	providerData, err := marshalProviderData(intent)
	if err != nil {
		return nil, err
	}

	sqlcIntent, err := repo.queries(ctx).UpdateIntentProviderData(ctx, sqlc.UpdateIntentProviderDataParams{
		ID:           intent.ID,
		ProviderData: providerData,
	})
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	return mapIntentFromDB(&sqlcIntent)
}

func (repo *intentsRepo) GetByMPOrderID(ctx context.Context, orderID string) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.GetByMPOrderID")
	defer span.End()

	sqlcIntent, err := repo.queries(ctx).GetIntentByMPOrderID(ctx, orderID)
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	return mapIntentFromDB(&sqlcIntent)
}

func (repo *intentsRepo) GetByMPTransactionID(ctx context.Context, transactionID string) (*models.Intent, error) {
	ctx, span := repo.tracer.Start(ctx, "IntentRepo.GetByMPTransactionID")
	defer span.End()

	sqlcIntent, err := repo.queries(ctx).GetIntentByMPTransactionID(ctx, transactionID)
	if err != nil {
		return nil, errx.FromDB(err, "intent")
	}

	return mapIntentFromDB(&sqlcIntent)
}
