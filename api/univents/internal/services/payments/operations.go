// Package payments wires an event to its Payssage payment configuration:
// the event's connected seller (provider account) on the single platform
// wallet (env-configured, D6) that receives the event's money.
//
// Payssage is an implementation detail — the Operations seam talks to it
// through PayssageClient (satisfied by sdk/payssage) and event owners never
// see wallets or callback URLs.
package payments

import (
	"context"
	"univents/internal/authz"
	"univents/ports"

	"sdk/payssage"

	"github.com/google/uuid"
)

// PayssageClient is the service-to-service seam into Payssage. It is
// satisfied by *payssage.Client (sdk/go/Payssage) and faked in tests.
type PayssageClient interface {
	ConnectProvider(ctx context.Context, provider string, req payssage.ConnectProviderRequest) (string, error)
	ListWalletSellers(ctx context.Context, walletID uuid.UUID) ([]payssage.Seller, error)
}

type Operations struct {
	events   ports.EventRepo
	payssage PayssageClient
	authz    *authz.Service
	walletID uuid.UUID // the single platform wallet every event's seller lives on (D6)
}

func NewOperations(
	events ports.EventRepo,
	payssage PayssageClient,
	authz *authz.Service,
	walletID uuid.UUID,
) *Operations {
	return &Operations{
		events:   events,
		payssage: payssage,
		authz:    authz,
		walletID: walletID,
	}
}

// ConnectResult is what the connect endpoint hands back to the caller: the
// provider consent URL to send the user to, plus the event's wallet id.
type ConnectResult struct {
	AuthURL  string    `json:"auth_url"`
	WalletID uuid.UUID `json:"wallet_id"`
}
