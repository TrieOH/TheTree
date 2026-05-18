package queries

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/google/uuid"
)

func (s *QueryService) BulkGet(ctx context.Context, ids []uuid.UUID) (ns []models.Namespace, err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.BulkGet")
	defer span.End()

	_, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	return s.namespaces.BulkGet(ctx, ids)
}
