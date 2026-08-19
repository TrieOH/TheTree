package intents

import (
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.IntentRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("intent"),
	}
}

func mapIntent(src sqlc.Intent) models.Intent {
	var statusDetail *models.IntentStatusDetail
	if src.StatusDetail != nil {
		statusDetail = new(models.IntentStatusDetail(*src.StatusDetail))
	}

	return models.Intent{
		ID:            src.ID,
		WalletID:      src.WalletID,
		SellerID:      src.SellerID,
		CollectorID:   src.CollectorID,
		AmountCents:   src.AmountCents,
		Currency:      src.Currency,
		Sandbox:       src.Sandbox,
		Provider:      src.Provider,
		Status:        models.IntentStatus(src.Status),
		StatusDetail:  statusDetail,
		ProviderData:  src.ProviderData,
		Metadata:      src.Metadata,
		ExternalID:    src.ExternalID,
		ExternalGroup: src.ExternalGroup,
		CreatedAt:     src.CreatedAt,
		UpdatedAt:     src.UpdatedAt,
	}
}
