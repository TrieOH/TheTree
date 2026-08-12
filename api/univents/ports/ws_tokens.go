package ports

import (
	"context"
	"univents/models"
)

// WsTokenRepo is the one-time handshake-token persistence (split 6). The
// token is stored hashed (SHA-256); only the hash is ever written or read.
// Create is transaction-aware through the shared tx runner (database.Queries
// picks the tx up from context), so checkout (split 7) inserts the token
// inside the checkout tx with the same call — no separate in-tx variant.
type WsTokenRepo interface {
	// Create inserts a ws_tokens row. The caller generates the raw token and
	// stores only its hash; the raw value is returned exactly once by the
	// issuer and never persisted.
	Create(ctx context.Context, toCreate *models.WsToken) (*models.WsToken, error)
	// Consume atomically marks a token used — WHERE token_hash matches AND
	// used_at IS NULL AND expires_at > now(). Returns (nil, nil) when the
	// guard misses (missing, already used, or expired): the handshake
	// rejects. One-time by construction.
	Consume(ctx context.Context, tokenHash string) (*models.WsToken, error)
}
