package commands

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (event *domain.Event, err error) {
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

	allowed, err := uc.gaClient.Authz.Check().User(sub.ID).
		Object("events").
		Action("edit").
		Scope(event.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("event").SetMessage("insufficient permissions")
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
