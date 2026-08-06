package badges

import (
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

// RevokeForRegistration revokes a participant badge when its registration is
// cancelled or expires. Called by the checkout feature.
func (o *Operations) RevokeForRegistration(ctx context.Context, registrationID uuid.UUID, reason string) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.RevokeForRegistration")
	defer span.End()
	return o.emissions.RevokeByRegistration(ctx, registrationID, reason)
}
