package infrastructure

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type eventsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ domain.EventsRepository = (*eventsRepo)(nil)

func NewEventRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) domain.EventsRepository {
	return &eventsRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *eventsRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapEventFromDB(src *sqlc.Event) *domain.Event {
	return &domain.Event{
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

func (repo *eventsRepo) CreateEvent(ctx context.Context, toCreate *domain.Event) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Create")
	defer span.End()

	event, err := repo.queries(ctx).CreateEvent(ctx, sqlc.CreateEventParams{
		ID:             toCreate.ID,
		OwnerID:        toCreate.OwnerID,
		OrganizationID: toCreate.OrganizationID,
		Name:           toCreate.Name,
		Acronym:        toCreate.Acronym,
		Slug:           toCreate.Slug,
		Tagline:        toCreate.Tagline,
		Description:    toCreate.Description,
		IsSeries:       toCreate.IsSeries,
		LogoUrl:        toCreate.LogoUrl,
		ContactEmail:   toCreate.ContactEmail,
		SocialLinks:    toCreate.SocialLinks,
		CreatedBy:      toCreate.CreatedBy,
		GoauthScopeID:  toCreate.GoauthScopeID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	return mapEventFromDB(&event), nil
}

func (repo *eventsRepo) PatchEvent(ctx context.Context, toPatch *domain.Event) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.PatchEvent")
	defer span.End()

	event, err := repo.queries(ctx).PatchEvent(ctx, sqlc.PatchEventParams{
		Name:          toPatch.Name,
		Acronym:       toPatch.Acronym,
		Slug:          toPatch.Slug,
		Tagline:       toPatch.Tagline,
		Description:   toPatch.Description,
		IsSeries:      toPatch.IsSeries,
		LogoUrl:       toPatch.LogoUrl,
		BannerUrl:     toPatch.BannerUrl,
		HasGallery:    toPatch.HasGallery,
		ContactEmail:  toPatch.ContactEmail,
		SocialLinks:   toPatch.SocialLinks,
		ID:            toPatch.ID,
		GoauthScopeID: toPatch.GoauthScopeID,
	})
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	return mapEventFromDB(&event), nil
}

func (repo *eventsRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.GetByID")
	defer span.End()

	event, err := repo.queries(ctx).GetEventByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	return mapEventFromDB(&event), nil
}

func (repo *eventsRepo) ListEvents(ctx context.Context) ([]domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Create")
	defer span.End()

	sqlcEvents, err := repo.queries(ctx).ListEvents(ctx)
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	outEvents := make([]domain.Event, 0, len(sqlcEvents))
	for _, sqlcEvent := range sqlcEvents {
		outEvents = append(outEvents, *mapEventFromDB(&sqlcEvent))
	}
	return outEvents, nil
}

func (repo *eventsRepo) ListOwnEvents(ctx context.Context, ownerID uuid.UUID) ([]domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.ListOwnEvents")
	defer span.End()

	sqlcEvents, err := repo.queries(ctx).ListOwnEvents(ctx, &ownerID)
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	outEvents := make([]domain.Event, 0, len(sqlcEvents))
	for _, sqlcEvent := range sqlcEvents {
		outEvents = append(outEvents, *mapEventFromDB(&sqlcEvent))
	}
	return outEvents, nil
}

func (repo *eventsRepo) PublishEvent(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.PublishEvent")
	defer span.End()

	err := repo.queries(ctx).PublishEvent(ctx, id)
	if err != nil {
		return errx.FromDB(err, "event")
	}

	return nil
}

func (repo *eventsRepo) AddEdition(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.AddEdition")
	defer span.End()

	affectedRows, err := repo.queries(ctx).AddEdition(ctx, eventID)
	if err != nil {
		return errx.FromDB(err, "event")
	}

	if affectedRows == 0 {
		return errx.Invalid("event").SetMessage("could not add edition")
	}

	return nil
}

func (repo *eventsRepo) AddGalleryImage(ctx context.Context, id uuid.UUID, url string) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.AddGalleryImage")
	defer span.End()

	sqlcEvent, err := repo.queries(ctx).AddEventGalleryImage(ctx, sqlc.AddEventGalleryImageParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	return mapEventFromDB(&sqlcEvent), nil
}

func (repo *eventsRepo) RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.RemoveGalleryImage")
	defer span.End()

	sqlcEvent, err := repo.queries(ctx).RemoveEventGalleryImage(ctx, sqlc.RemoveEventGalleryImageParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	return mapEventFromDB(&sqlcEvent), nil
}

func (repo *eventsRepo) SetLogo(ctx context.Context, id uuid.UUID, url string) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.SetLogo")
	defer span.End()

	sqlcEvent, err := repo.queries(ctx).SetEventLogo(ctx, sqlc.SetEventLogoParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	return mapEventFromDB(&sqlcEvent), nil
}

func (repo *eventsRepo) UnsetLogo(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.UnsetLogo")
	defer span.End()

	sqlcEvent, err := repo.queries(ctx).UnsetEventLogo(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	return mapEventFromDB(&sqlcEvent), nil
}

func (repo *eventsRepo) SetBanner(ctx context.Context, id uuid.UUID, url string) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.SetBanner")
	defer span.End()

	sqlcEvent, err := repo.queries(ctx).SetEventBanner(ctx, sqlc.SetEventBannerParams{
		ID:  id,
		Url: url,
	})
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	return mapEventFromDB(&sqlcEvent), nil
}

func (repo *eventsRepo) UnsetBanner(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.UnsetBanner")
	defer span.End()

	sqlcEvent, err := repo.queries(ctx).UnsetEventBanner(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "event")
	}

	return mapEventFromDB(&sqlcEvent), nil
}
