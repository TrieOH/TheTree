package ports

import (
	"context"
	"payssage/models"
)

type CollectorRepo interface {
	Create(ctx context.Context, toCreate models.Collector) (collector models.Collector, err error)
}
