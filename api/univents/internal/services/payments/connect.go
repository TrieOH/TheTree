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

// Connect ensures the event's wallet exists (creating it lazily under the
// platform identity, with the 5% marketplace fee) and starts the provider
// OAuth flow for a seller on that wallet. The callback URL is assembled per
// event, at call time: the event id goes in the path, so Payssage's own
// `?credential_id=…&public_key=…` append keeps working untouched.
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

	walletID := event.PayssageWalletID
	if walletID == nil {
		wallet, err := o.payssage.CreateWallet(ctx, payssage.CreateWalletRequest{Name: event.Slug})
		if err != nil {
			return nil, err
		}
		err = o.payssage.SetWalletFee(ctx, wallet.ID, eventWalletFeeBps)
		if err != nil {
			return nil, err
		}
		// Persist the wallet on the event before starting the OAuth flow so a
		// retry reuses it instead of orphaning another wallet.
		event, err = o.events.SetPaymentsConfig(ctx, event.ID, nil, &wallet.ID, nil)
		if err != nil {
			return nil, err
		}
		walletID = &wallet.ID
	}

	finalRedirectURL := fmt.Sprintf("%s/events/%s/payssage/oauth/callback",
		strings.TrimRight(os.Getenv("APP_URL"), "/"), event.ID)

	authURL, err := o.payssage.ConnectProvider(ctx, provider, payssage.ConnectProviderRequest{
		Flow:             payssage.OAuthFlowSeller,
		WalletID:         walletID,
		FinalRedirectURL: finalRedirectURL,
	})
	if err != nil {
		return nil, err
	}

	return &ConnectResult{AuthURL: authURL, WalletID: *walletID}, nil
}
