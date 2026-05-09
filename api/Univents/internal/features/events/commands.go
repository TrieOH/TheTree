package events

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"
	"univents/internal/platform/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"

	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	events ports.EventsRepository
	minio  *minio.Client
	tracer trace.Tracer
	az     *authzed.Client
	tx     database.TxRunner
}

func NewCommandService(
	events ports.EventsRepository,
	minio *minio.Client,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		events: events,
		minio:  minio,
		tracer: tracer,
		az:     az,
		tx:     tx,
	}
}

func (uc *CommandService) CreateEvent(ctx context.Context, in contracts.CreateEventSpec) (out *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var validEvent *contracts.Event
	validEvent, err = contracts.NewEvent(sub.ID, &sub.ID, in)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_events"),
		authz.Resource("platform", "global"),
	); err != nil {
		return nil, err
	}

	if err = authz.GrantRole(ctx, uc.az, "platform:global#event_creator@user:"+sub.ID.String()); err != nil {
		return nil, err
	} // FIXME Outbox this too

	var created *contracts.Event
	created, err = uc.events.CreateEvent(ctx, validEvent) // FIXME if this fails the scope must be undone (SAGA PATTERN)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *CommandService) PublishEvent(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "EventService.PublishEvent")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	event, err := uc.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("publish"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
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

func (uc *CommandService) PatchEvent(ctx context.Context, in contracts.PatchEventSpec) (out *contracts.Event, warns []string, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.PatchEvent")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("patch.success", err == nil))
	}()

	span.SetAttributes(attribute.String("event.id", in.ID.String()))

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, nil, err
	}

	var event *contracts.Event
	event, err = uc.events.GetByID(ctx, in.ID)
	if err != nil {
		return nil, nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
		return nil, nil, err
	}

	if event.Name != in.Name {
		event.Name = in.Name
	}

	if in.Acronym != nil && (event.Acronym == nil || *in.Acronym != *event.Acronym) {
		event.Acronym = in.Acronym
	} else if in.Acronym == nil && event.Acronym != nil {
		event.Acronym = in.Acronym
	}

	if event.Slug != in.Slug {
		event.Slug = in.Slug
	}

	if in.Tagline != nil && (event.Tagline == nil || *in.Tagline != *event.Tagline) {
		event.Tagline = in.Tagline
	} else if in.Tagline == nil && event.Tagline != nil {
		event.Tagline = in.Tagline
	}

	if in.Description != nil && (event.Description == nil || *in.Description != *event.Description) {
		event.Description = in.Description
	} else if in.Description == nil && event.Description != nil {
		event.Description = in.Description
	}

	if event.IsSeries != in.IsSeries {
		if !in.IsSeries && event.EditionsCount > 1 {
			warns = append(warns, "Cannot convert to non-series when multiple editions exist")
		} else {
			event.IsSeries = in.IsSeries
		}
	}

	if in.ContactEmail != nil && (event.ContactEmail == nil || *in.ContactEmail != *event.ContactEmail) {
		event.ContactEmail = in.ContactEmail
	} else if in.ContactEmail == nil && event.ContactEmail != nil {
		event.ContactEmail = in.ContactEmail
	}

	if in.SocialLinks != nil && (event.SocialLinks == nil || !bytes.Equal(*in.SocialLinks, *event.SocialLinks)) {
		event.SocialLinks = in.SocialLinks
	} else if in.SocialLinks == nil && event.SocialLinks != nil {
		event.SocialLinks = in.SocialLinks
	}

	var patched *contracts.Event
	patched, err = uc.events.PatchEvent(ctx, event)
	if err != nil {
		return nil, warns, err
	}

	return patched, warns, nil
}

func (uc *CommandService) SetLogo(ctx context.Context, id uuid.UUID, url string) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.SetLogo")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("set_logo.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
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

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event.LogoUrl == nil {
		return nil, errx.Invalid("event").SetMessage("already has no logo")
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
		return nil, err
	}

	bucket, key, err := parseMinioURL(*event.LogoUrl)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.minio.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{}); err != nil {
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

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
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

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event.BannerUrl == nil {
		return nil, errx.Invalid("event").SetMessage("already has no banner")
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
		return nil, err
	}

	bucket, key, err := parseMinioURL(*event.BannerUrl)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.minio.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{}); err != nil {
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

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
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

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
		return nil, err
	}

	bucket, key, err := parseMinioURL(url)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.minio.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return nil, errx.Internal("event").SetMessage("failed to delete image from storage: " + err.Error())
	}

	event, err = uc.events.RemoveGalleryImage(ctx, event.ID, url)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func parseMinioURL(rawURL string) (bucket, key string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid url: %w", err)
	}
	// path is /bucket/key/possibly/nested
	parts := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("url path too short, expected /bucket/key: %s", u.Path)
	}
	return parts[0], parts[1], nil
}
