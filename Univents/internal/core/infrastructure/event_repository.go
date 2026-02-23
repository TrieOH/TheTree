package infrastructure

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
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

func EventStatusPtrToString(p *domain.EventStatus) string {
	if p == nil {
		return ""
	}
	return string(*p)
}

func StringPtrToEventStatusPtr(p *string) *domain.EventStatus {
	var stat domain.EventStatus
	if p == nil {
		return nil
	}
	stat = domain.EventStatus(*p)
	return &stat
}

func mapEventAuditFromDB(src *sqlc.EventAudit) *domain.Audit {
	from := EventStatusPtrToString(src.FromStatus)
	to := EventStatusPtrToString(src.ToStatus)
	return &domain.Audit{
		ID:         src.ID,
		ResourceID: src.EventID,
		ActorType:  domain.ActorType(src.ActorType),
		ActorID:    src.ActorID,
		Action:     string(src.Action),
		State:      domain.AuditActionState(src.State),
		FromStatus: &from,
		ToStatus:   &to,
		Metadata:   src.Metadata,
		CreatedAt:  src.CreatedAt,
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
		return nil, fail.From(err).RecordCtx(ctx)
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
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return mapEventFromDB(&event), nil
}

func (repo *eventsRepo) GetEventByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.GetEventByID")
	defer span.End()

	event, err := repo.queries(ctx).GetEventByID(ctx, id)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return mapEventFromDB(&event), nil
}

func (repo *eventsRepo) ListEvents(ctx context.Context) ([]domain.Event, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.Create")
	defer span.End()

	sqlcEvents, err := repo.queries(ctx).ListEvents(ctx)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
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
		return nil, fail.From(err).RecordCtx(ctx)
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
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *eventsRepo) AppendEventAudits(ctx context.Context, audits []domain.Audit) error {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.AppendEventAudits")
	defer span.End()

	params := make([]sqlc.AppendEventAuditsParams, len(audits))
	for i, a := range audits {
		params[i] = sqlc.AppendEventAuditsParams{
			EventID:    a.ResourceID,
			ActorType:  sqlc.AuditActorType(a.ActorType),
			ActorID:    a.ActorID,
			Action:     sqlc.EventAuditAction(a.Action),
			FromStatus: StringPtrToEventStatusPtr(a.FromStatus),
			ToStatus:   StringPtrToEventStatusPtr(a.ToStatus),
			State:      sqlc.ActionState(a.State),
			Metadata:   a.Metadata,
		}
	}

	amount, err := repo.queries(ctx).AppendEventAudits(ctx, params)
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Int64("created.amount", amount))
	return nil
}

func (repo *eventsRepo) ListEventAuditByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Audit, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.ListEventAuditByEvent")
	defer span.End()

	sqlcAudits, err := repo.queries(ctx).ListEventAuditByEvent(ctx, eventID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	outAudits := make([]domain.Audit, 0, len(sqlcAudits))
	for _, audit := range sqlcAudits {
		outAudits = append(outAudits, *mapEventAuditFromDB(&audit))
	}
	return outAudits, nil
}

func (repo *eventsRepo) AddEdition(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.AddEdition")
	defer span.End()

	affectedRows, err := repo.queries(ctx).AddEdition(ctx, eventID)
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	if affectedRows == 0 {
		return fail.New(errx.EventCannotAddEditions).RecordCtx(ctx)
	}

	return nil
}
