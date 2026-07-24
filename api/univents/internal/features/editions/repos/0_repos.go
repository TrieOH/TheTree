package repos

import (
	"univents/internal/database/sqlc"
	"univents/models"
	"univents/ports"

	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.EditionRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.EditionRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("edition"),
	}
}

func mapEdition(src sqlc.Edition) models.Edition {
	return models.Edition{
		ID:                  src.ID,
		EventID:             src.EventID,
		Name:                src.EditionName,
		Slug:                src.Slug,
		Tagline:             src.Tagline,
		Description:         src.Description,
		IsDraft:             src.IsDraft,
		RegistrationOpensAt: src.RegistrationOpensAt,
		StartsAt:            src.StartsAt,
		EndsAt:              src.EndsAt,
		LocationName:        src.LocationName,
		LocationAddress:     src.LocationAddress,
		LogoURL:             src.LogoUrl,
		BannerURL:           src.BannerUrl,
		ContactEmail:        src.ContactEmail,
		CreatedBy:           src.CreatedBy,
		CreatedAt:           src.CreatedAt,
		UpdatedAt:           src.UpdatedAt,
		DeletedAt:           src.DeletedAt,
	}
}
