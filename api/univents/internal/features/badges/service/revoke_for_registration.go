package service

import (
	"context"

	"github.com/google/uuid"
)

// RevokeForRegistration revokes a participant badge when its registration is
// cancelled or expires. Called by the checkout feature.
func (o *Operations) RevokeForRegistration(ctx context.Context, registrationID uuid.UUID, reason string) error {
	return o.Emissions.RevokeByRegistration(ctx, registrationID, reason)
}
