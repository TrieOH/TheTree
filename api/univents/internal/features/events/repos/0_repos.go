package repos

import (
	"univents/internal/database/sqlc"
	"univents/models"

	"lib/database"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.EventRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.EventRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("event"),
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
