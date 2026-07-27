package repos

import (
	sqlc2 "univents/internal/sqlc"
	"univents/models"

	"lib/database"
	"univents/ports"
)

type Repo struct {
	q   *sqlc2.Queries
	dbe database.ErrorHandler
}

var _ ports.EventRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("event"),
	}
}

func mapEventMember(src sqlc2.EventMember) models.EventMember {
	return models.EventMember{
		ID:        src.ID,
		EventID:   src.EventID,
		UserID:    src.UserID,
		Role:      models.EventMemberRole(src.Role),
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		DeletedAt: src.DeletedAt,
	}
}

func mapEvent(src sqlc2.Event) models.Event {
	return models.Event{
		ID:               src.ID,
		OwnerID:          src.OwnerID,
		FullName:         src.FullName,
		Acronym:          src.Acronym,
		Slug:             src.Slug,
		Description:      src.Description,
		Style:            src.Style,
		Status:           models.EventStatus(src.Status),
		PayssageSellerID: src.PayssageSellerID,
		PayssageWalletID: src.PayssageWalletID,
		LogoURL:          src.LogoUrl,
		BannerURL:        src.BannerUrl,
		ContactEmail:     src.ContactEmail,
		CreatedAt:        src.CreatedAt,
		UpdatedAt:        src.UpdatedAt,
		DeletedAt:        src.DeletedAt,
	}
}
