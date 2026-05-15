package app

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
	"univents/internal/platform/database/sqlc"
	"univents/internal/platform/telemetry"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func productHardDeleteJob(scheduler gocron.Scheduler, objectStorage *minio.Client, db *pgxpool.Pool) {
	_, err := scheduler.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(func(ctx context.Context, pool *pgxpool.Pool) {
			ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()

			q := sqlc.New(pool)

			products, err := q.GetExpiredSoftDeletedProducts(ctx)
			if err != nil {
				telemetry.Log().Error("Couldn't fetch expired soft-deleted products", zap.Error(err))
				return
			}

			if len(products) == 0 {
				return
			}

			var hardDeleted []uuid.UUID
			for _, product := range products {
				failed := false

				urls := collectURLs(product.ThumbnailUrl, product.GalleryUrls)
				for _, u := range urls {
					bucket, key, err := parseMinioURL(u)
					if err != nil {
						telemetry.Log().Error("Couldn't parse MinIO URL, skipping",
							zap.String("url", u),
							zap.Error(err),
						)
						failed = true
						continue
					}
					if err := objectStorage.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{}); err != nil {
						telemetry.Log().Error("Couldn't delete object from MinIO, skipping product",
							zap.String("bucket", bucket),
							zap.String("key", key),
							zap.Error(err),
						)
						failed = true
						break
					}
				}

				if !failed {
					hardDeleted = append(hardDeleted, product.ID)
				}
			}

			if len(hardDeleted) == 0 {
				return
			}

			if err := q.MarkProductsHardDeleted(ctx, hardDeleted); err != nil {
				telemetry.Log().Error("Couldn't mark products as hard deleted", zap.Error(err))
				return
			}

			telemetry.Log().Info("Hard deleted expired products", zap.Int("count", len(hardDeleted)))
		}, db),
	)

	if err != nil {
		log.Fatalf("Couldn't create ProductHardDelete cron job: %v", err)
	}
	log.Println("Created ProductHardDelete cron job")
}

// collectURLs flattens thumbnail + gallery into a single slice, skipping nulls.
func collectURLs(thumbnail *string, gallery []string) []string {
	var urls []string
	if thumbnail != nil {
		urls = append(urls, *thumbnail)
	}
	urls = append(urls, gallery...)
	return urls
}

// parseMinioURL extracts bucket and key from http://host:port/bucket/key/path
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
