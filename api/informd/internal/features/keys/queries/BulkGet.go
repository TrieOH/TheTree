package queries

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/google/uuid"
)

func (s *QueryService) BulkGet(ctx context.Context, ids []uuid.UUID) (keys []models.APIKey, err error) {
	ctx, span := s.tracer.Start(ctx, "ApiKeyService.BulkGet")
	defer span.End()

	_, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	return s.apiKeys.BulkGet(ctx, ids)
}
