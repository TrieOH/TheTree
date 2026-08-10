package payments

import (
	"context"
	"fmt"
	"lib/telemetry"
	"os"
	"strings"
	"univents/models"

	idx "sdk/identityx"
	payssage "sdk/payssage"

	"github.com/google/uuid"
)

// Connect starts the provider OAuth flow for a seller on the platform
// wallet (env-configured, D6) — no per-event wallet exists anymore. The
// callback URL is assembled per event, at call time: the event id goes in
// the path, so Payssage's own `?credential_id=…&public_key=…` append keeps
// working untouched.
func (o *Operations) Connect(ctx context.Context, eventID uuid.UUID, provider string) (*ConnectResult, error) {
	ctx, span := telemetry.StartSpan(ctx, "PaymentsService.Connect")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	finalRedirectURL := fmt.Sprintf("%s/events/%s/payssage/oauth/callback",
		strings.TrimRight(os.Getenv("APP_URL"), "/"), event.ID)

	authURL, err := o.payssage.ConnectProvider(ctx, provider, payssage.ConnectProviderRequest{
		Flow:             payssage.OAuthFlowSeller,
		WalletID:         &o.walletID,
		FinalRedirectURL: finalRedirectURL,
	})
	if err != nil {
		return nil, err
	}

	return &ConnectResult{AuthURL: authURL, WalletID: o.walletID}, nil
}
