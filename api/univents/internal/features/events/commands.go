package events

import (
	"context"
	"lib/database"
	"lib/objectstorage"
	"univents/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"

	idx "sdk/identityx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	events ports.EventsRepository
	obj    *objectstorage.Client
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewCommandService(
	events ports.EventsRepository,
	obj *objectstorage.Client,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		events: events,
		obj:    obj,
		logger: logger,
		tracer: tracer,
		tx:     tx,
	}
}

func (uc *CommandService) CreateEvent(ctx context.Context, in contracts.CreateEventSpec) (out *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var validEvent *contracts.Event
	validEvent, err = contracts.NewEvent(ident.Sub.ID, &ident.Sub.ID, in)
	if err != nil {
		return nil, err
	}

	var created *contracts.Event
	created, err = uc.events.CreateEvent(ctx, validEvent)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *CommandService) PublishEvent(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "EventService.PublishEvent")
	defer span.End()

	event, err := uc.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != contracts.StatusDraft {
		return errx.Invalid("event").SetMessage("cannot publish non draft event")
	}

	err = uc.events.PublishEvent(ctx, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) SetLogo(ctx context.Context, id uuid.UUID, url string) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.SetLogo")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("set_logo.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.SetLogo(ctx, event.ID, url)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (uc *CommandService) UnsetLogo(ctx context.Context, id uuid.UUID) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.UnsetLogo")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("unset_logo.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event.LogoUrl == nil {
		return nil, errx.Invalid("event").SetMessage("already has no logo")
	}

	bucket, key, err := objectstorage.ParseURL(*event.LogoUrl)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.obj.RemoveObject(ctx, bucket, key); err != nil {
		return nil, errx.Internal("event").SetMessage("failed to delete image from storage: " + err.Error())
	}

	event, err = uc.events.UnsetLogo(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (uc *CommandService) SetBanner(ctx context.Context, id uuid.UUID, url string) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.SetBanner")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("set_banner.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.SetBanner(ctx, event.ID, url)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (uc *CommandService) UnsetBanner(ctx context.Context, id uuid.UUID) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.UnsetBanner")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("unset_banner.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event.BannerUrl == nil {
		return nil, errx.Invalid("event").SetMessage("already has no banner")
	}

	bucket, key, err := objectstorage.ParseURL(*event.BannerUrl)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.obj.RemoveObject(ctx, bucket, key); err != nil {
		return nil, errx.Internal("event").SetMessage("failed to delete image from storage: " + err.Error())
	}

	event, err = uc.events.UnsetBanner(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (uc *CommandService) AddGalleryImage(ctx context.Context, id uuid.UUID, url string) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.AddGalleryImage")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("add_gallery.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.AddGalleryImage(ctx, event.ID, url)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (uc *CommandService) RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.RemoveGalleryImage")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("remove_gallery.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	bucket, key, err := objectstorage.ParseURL(url)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.obj.RemoveObject(ctx, bucket, key); err != nil {
		return nil, errx.Internal("event").SetMessage("failed to delete image from storage: " + err.Error())
	}

	event, err = uc.events.RemoveGalleryImage(ctx, event.ID, url)
	if err != nil {
		return nil, err
	}

	return event, nil
}
