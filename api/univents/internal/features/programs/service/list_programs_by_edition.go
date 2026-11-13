package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListProgramsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Program, error) {
	return o.Programs.ListByEdition(ctx, editionID)
}
