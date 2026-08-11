package purchases

import (
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.PurchaseRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("purchase"),
	}
}

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
