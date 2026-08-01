package signatures

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.ListByEdition")
	defer span.End()
	return o.signatures.ListByEdition(ctx, editionID)
}
