package repos

import (
	"univents/internal/database/sqlc"

	"lib/database"
	"univents/contracts"
	"univents/internal/shared/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.EventsRepository = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.EventsRepository {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("event"),
	}
}

func mapEventFromDB(src *sqlc.Event) *contracts.Event {
	return &contracts.Event{
		ID:             src.ID,
		OwnerID:        src.OwnerID,
		OrganizationID: src.OrganizationID,
		GoauthScopeID:  src.GoauthScopeID,
		Name:           src.Name,
		Acronym:        src.Acronym,
		Slug:           src.Slug,
		Tagline:        src.Tagline,
		Description:    src.Description,
		IsSeries:       src.IsSeries,
		EditionsCount:  src.EditionsCount,
		LogoUrl:        src.LogoUrl,
		BannerUrl:      src.BannerUrl,
		HasGallery:     src.HasGallery,
		GalleryUrls:    src.GalleryUrls,
		ContactEmail:   src.ContactEmail,
		SocialLinks:    src.SocialLinks,
		Status:         src.Status,
		CreatedBy:      src.CreatedBy,
		CreatedAt:      src.CreatedAt,
		UpdatedAt:      src.UpdatedAt,
		DeletedAt:      src.DeletedAt,
	}
}
