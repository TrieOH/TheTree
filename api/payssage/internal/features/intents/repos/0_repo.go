package repos

import (
	"lib/database"
	sqlc2 "payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
)

type repo struct {
	q      *sqlc2.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.IntentRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.IntentRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("intent"),
	}
}

func mapIntent(src sqlc2.Intent) models.Intent {
	var statusDetail *models.IntentStatusDetail
	if src.StatusDetail != nil {
		statusDetail = new(models.IntentStatusDetail(*src.StatusDetail))
	}

	return models.Intent{
		ID:           src.ID,
		WalletID:     src.WalletID,
		SellerID:     src.SellerID,
		CollectorID:  src.CollectorID,
		AmountCents:  src.AmountCents,
		Currency:     src.Currency,
		Sandbox:      src.Sandbox,
		Provider:     src.Provider,
		Status:       models.IntentStatus(src.Status),
		StatusDetail: statusDetail,
		ProviderData: src.ProviderData,
		Metadata:     src.Metadata,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}
