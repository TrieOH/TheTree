package queries

import (
	"context"
	"lib/crypto"

	"go.uber.org/zap"
)

func (c *Queries) JWKS(ctx context.Context) (map[string]any, error) {
	keys, err := c.cryptoKeys.GetActiveSigningKeys(ctx, nil)
	if err != nil {
		return nil, err
	}

	jwks := make([]map[string]any, 0, len(keys))
	for _, k := range keys {
		jwk, err := crypto.PublicKeyToJWKS(k.ID.String(), k.PublicKey)
		if err != nil {
			c.logger.Warn("skipping malformed key", zap.String("key_id", k.ID.String()), zap.Error(err))
			continue
		}
		jwks = append(jwks, jwk)
	}

	return map[string]any{"keys": jwks}, nil
}
