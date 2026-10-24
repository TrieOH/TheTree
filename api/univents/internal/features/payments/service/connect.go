package service

import (
	"context"
	"fmt"
	"os"
	payssage "sdk/payssage"
	"strings"
	"univents/models"

	"github.com/google/uuid"
)

// Connect starts the provider OAuth flow for a seller on the platform
// wallet (env-configured, D6) — no per-event wallet exists anymore. The
// callback URL is assembled per event, at call time: the event id goes in
// the path, so Payssage's own `?credential_id=…&public_key=…` append keeps
// working untouched.
func (o *Operations) Connect(ctx context.Context, eventID uuid.UUID, provider string) (*ConnectResult, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	finalRedirectURL := fmt.Sprintf("%s/events/%s/payssage/oauth/callback",
		strings.TrimRight(os.Getenv("APP_URL"), "/"), event.ID)

	authURL, err := o.Payssage.ConnectProvider(ctx, provider, payssage.ConnectProviderRequest{
		Flow:             payssage.OAuthFlowSeller,
		WalletID:         &o.WalletID,
		FinalRedirectURL: finalRedirectURL,
	})
	if err != nil {
		return nil, err
	}

	return &ConnectResult{AuthURL: authURL, WalletID: o.WalletID}, nil
}
