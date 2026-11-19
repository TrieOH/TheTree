// Package payments is the payments feature: its handlers and operations.
// Build wires it.
//
// The operations wire an event to its Payssage payment configuration:
// the event's connected seller (provider account) on the single platform
// wallet (env-configured, D6) that receives the event's money.
//
// Payssage is an implementation detail — the Operations seam talks to it
// through PayssageClient (satisfied by sdk/payssage) and event owners never
// see wallets or callback URLs.
package payments

import (
	"context"
	"fmt"
	payssage "sdk/payssage"
	"strings"
	"univents/errs"
	"univents/internal/authz"
	"univents/models"
	"univents/ports"
	"uuid"
)

// PayssageClient is the service-to-service seam into Payssage. It is
// satisfied by *payssage.Client (sdk/go/Payssage) and faked in tests.
type PayssageClient interface {
	ConnectProvider(ctx context.Context, provider string, req payssage.ConnectProviderRequest) (string, error)
	ListWalletSellers(ctx context.Context, walletID uuid.UUID) ([]payssage.Seller, error)
}

type Deps struct {
	Events   ports.EventRepo
	Payssage PayssageClient
	Authz    *authz.Service
	WalletID uuid.UUID // the single platform wallet every event's seller lives on (D6)
	AppURL   string    // the Univents front, where the OAuth flow lands back
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// ConnectResult is what the connect endpoint hands back to the caller: the
// provider consent URL to send the user to, plus the event's wallet id.
type ConnectResult struct {
	AuthURL  string    `json:"auth_url"`
	WalletID uuid.UUID `json:"wallet_id"`
}

// Complete persists the seller + provider public key on the event, after the
// OAuth callback delivered them (seller_id = payssage's credential_id). The
// seller must belong to the platform wallet (env-configured, D6) — verified
// against Payssage, so a forged or foreign seller id is rejected.
func (o *Operations) Complete(ctx context.Context, eventID, sellerID uuid.UUID, publicKey string) (*models.Event, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	sellers, err := o.Payssage.ListWalletSellers(ctx, o.WalletID)
	if err != nil {
		return nil, err
	}
	belongs := false
	for _, s := range sellers {
		if s.ID == sellerID {
			belongs = true
			break
		}
	}
	if !belongs {
		return nil, errs.ErrSellerNotInPlatformWallet
	}

	return o.Events.SetPaymentsConfig(ctx, event.ID, &sellerID, &publicKey)
}

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
		strings.TrimRight(o.AppURL, "/"), event.ID)

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

// Disconnect unlinks the seller from the event (unlink only): the seller and
// the provider public key are cleared. The platform wallet itself is
// env-configured (D6) and untouched — reconnecting just re-OAuths a seller
// onto it.
func (o *Operations) Disconnect(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Events.ClearPaymentsConfig(ctx, event.ID)
}
