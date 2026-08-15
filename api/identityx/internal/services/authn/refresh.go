package authn

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"
)

// Refresh rotates the session: the tokens module verifies the old refresh
// token, blacklists it and the access token it anchors, and mints a fresh
// pair for the actor. Fail-closed — if the old pair cannot be revoked, no
// new pair is issued.
func (o *Operations) Refresh(ctx context.Context, refreshToken string) (*models.UserTokensOutput, error) {
	ctx, span := telemetry.StartSpan(ctx, "Refresh")
	defer span.End()

	return o.tokens.Rotate(ctx, refreshToken)
}
