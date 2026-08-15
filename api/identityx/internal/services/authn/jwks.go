package authn

import (
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

// JWKS publishes the project's public signing keys as a JWKS document. The
// tokens module owns key policy: scope validation and the key read live
// there, not here.
func (o *Operations) JWKS(ctx context.Context, projectID *uuid.UUID) (map[string]any, error) {
	ctx, span := telemetry.StartSpan(ctx, "JWKS")
	defer span.End()

	return o.tokens.JWKS(ctx, projectID)
}
