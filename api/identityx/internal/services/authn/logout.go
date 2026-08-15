package authn

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"
)

// Logout ends the session: the tokens module blacklists the access token
// and, when the refresh token is still verifiable, the refresh token too.
// A dead refresh token never fails the logout — the access token is
// already revoked, so the session is over either way.
func (o *Operations) Logout(ctx context.Context, in models.LogoutInput) error {
	ctx, span := telemetry.StartSpan(ctx, "Logout")
	defer span.End()

	return o.tokens.Revoke(ctx, in.AccessToken, in.RefreshToken)
}
