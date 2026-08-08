// Package payments wires an event to its Payssage payment configuration:
// the event's wallet (created lazily, owned by the platform identity) and the
// connected seller (provider account) that receives the event's money.
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

// eventWalletFeeBps is the 5% marketplace fee applied to every event wallet,
// set once at wallet creation.
const eventWalletFeeBps = 500

// PayssageClient is the service-to-service seam into Payssage. It is
// satisfied by *payssage.Client (sdk/go/Payssage) and faked in tests.
type PayssageClient interface {
	CreateWallet(ctx context.Context, req payssage.CreateWalletRequest) (*payssage.Wallet, error)
	SetWalletFee(ctx context.Context, walletID uuid.UUID, feeBps int) error
	ConnectProvider(ctx context.Context, provider string, req payssage.ConnectProviderRequest) (string, error)
	ListWalletSellers(ctx context.Context, walletID uuid.UUID) ([]payssage.Seller, error)
}

type Operations struct {
	events   ports.EventRepo
	payssage PayssageClient
	authz    *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	payssage PayssageClient,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:   events,
		payssage: payssage,
		authz:    authz,
	}
}

// ConnectResult is what the connect endpoint hands back to the caller: the
// provider consent URL to send the user to, plus the event's wallet id.
type ConnectResult struct {
	AuthURL  string    `json:"auth_url"`
	WalletID uuid.UUID `json:"wallet_id"`
}
