package infrastructure

import (
	"context"
	"univents/internal/eventcore/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"

	"github.com/MintzyG/fail/v3"
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

func mapEventFromDB(dst *domain.Event, src *sqlc.Event) {
	dst.ID = src.ID
	dst.OrganizationID = src.OrganizationID
	dst.Name = src.Name
	dst.Acronym = src.Acronym
	dst.Slug = src.Slug
	dst.Tagline = src.Tagline
	dst.Description = src.Description
	dst.IsSeries = src.IsSeries
	dst.EditionsCount = src.EditionsCount
	dst.LogoUrl = src.LogoUrl
	dst.BannerUrl = src.BannerUrl
	dst.HasGallery = src.HasGallery
	dst.GalleryUrls = src.GalleryUrls
	dst.ContactEmail = src.ContactEmail
	dst.SocialLinks = src.SocialLinks
	dst.Status = domain.Status(src.Status)
	dst.CreatedBy = src.CreatedBy
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.DeletedAt = src.DeletedAt
}

func mapEventAuditFromDB(dst *domain.Audit, src *sqlc.EventAudit) {
	dst.ID = src.ID
	dst.EventID = src.EventID
	dst.ActorType = domain.ActorType(src.ActorType)
	dst.ActorID = src.ActorID
	dst.Action = domain.AuditAction(src.Action)
	dst.Metadata = src.Metadata
	dst.CreatedAt = src.CreatedAt
	dst.FromStatus = src.FromStatus
	dst.ToStatus = src.ToStatus
}

func (repo *eventsRepo) Create(ctx context.Context, toCreate domain.Event) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Create")
	defer span.End()

	event, err := repo.queries(ctx).CreateEvent(ctx, sqlc.CreateEventParams{
		ID:             toCreate.ID,
		OrganizationID: toCreate.OrganizationID,
		Name:           toCreate.Name,
		Acronym:        toCreate.Acronym,
		Slug:           toCreate.Slug,
		Tagline:        toCreate.Tagline,
		IsSeries:       toCreate.IsSeries,
		LogoUrl:        toCreate.LogoUrl,
		ContactEmail:   toCreate.ContactEmail,
		SocialLinks:    toCreate.SocialLinks,
		CreatedBy:      toCreate.CreatedBy,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var out domain.Event
	mapEventFromDB(&out, &event)
	return &out, nil
}

func (repo *eventsRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.GetByID")
	defer span.End()

	event, err := repo.queries(ctx).GetEventByID(ctx, id)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var out domain.Event
	mapEventFromDB(&out, &event)
	return &out, nil
}

func (repo *eventsRepo) List(ctx context.Context) ([]domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Create")
	defer span.End()

	sqlcEvents, err := repo.queries(ctx).ListEvents(ctx)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	outEvents := make([]domain.Event, 0, len(sqlcEvents))
	for _, sqlcEvent := range sqlcEvents {
		var outEvent domain.Event
		mapEventFromDB(&outEvent, &sqlcEvent)
		outEvents = append(outEvents, outEvent)
	}
	return outEvents, nil
}

func (repo *eventsRepo) Publish(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Publish")
	defer span.End()

	err := repo.queries(ctx).PublishEvent(ctx, id)
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *eventsRepo) AppendAudit(ctx context.Context, audit domain.Audit) (*domain.Audit, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Publish")
	defer span.End()

	sqlcAudit, err := repo.queries(ctx).AppendEventAudit(ctx, sqlc.AppendEventAuditParams{
		EventID:    audit.EventID,
		ActorType:  sqlc.EventActorType(audit.ActorType),
		ActorID:    audit.ActorID,
		Action:     sqlc.EventAuditAction(audit.Action),
		FromStatus: audit.FromStatus,
		ToStatus:   audit.ToStatus,
		Metadata:   audit.Metadata,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var auditResult domain.Audit
	mapEventAuditFromDB(&auditResult, &sqlcAudit)
	return &auditResult, nil
}
