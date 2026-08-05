package signatures

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.GetByID")
	defer span.End()
	return o.signatures.GetByID(ctx, id)
}
