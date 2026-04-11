package commands

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (product *domain.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.RemoveGalleryImage")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("remove_gallery.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return nil, err
	}

	bucket, key, err := parseMinioURL(url)
	if err != nil {
		return nil, errx.Invalid("product").SetMessage("invalid image url")
	}

	if err = uc.minio.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return nil, errx.Internal("product").SetMessage("failed to delete image from storage: " + err.Error())
	}

	product, err = uc.products.RemoveGalleryImage(ctx, product.ID, url)
	if err != nil {
		return nil, err
	}

	return product, nil
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
