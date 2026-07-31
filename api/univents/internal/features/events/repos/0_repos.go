package repos

import (
	"univents/internal/sqlc"
	"univents/models"

	"lib/database"
	"univents/ports"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.EventRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("event"),
	}
}

func mapEventMember(src sqlc.EventMember) models.EventMember {
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

func mapEvent(src sqlc.Event) models.Event {
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
