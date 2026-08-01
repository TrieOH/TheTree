package authn

import (
	"context"
	"lib/crypto"
	"lib/telemetry"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (o *Operations) JWKS(ctx context.Context, projectID *uuid.UUID) (map[string]any, error) {
	ctx, span := telemetry.StartSpan(ctx, "JWKS")
	defer span.End()

	if projectID != nil {
		_, err := o.projects.GetByID(ctx, *projectID)
		if err != nil {
			return nil, err
		}
	}

	keys, err := o.cryptoKeys.GetActiveSigningKeys(ctx, projectID)
	if err != nil {
		return nil, err
	}

	jwks := make([]map[string]any, 0, len(keys))
	for _, k := range keys {
		jwk, err := crypto.PublicKeyToJWKS(k.ID.String(), k.PublicKey)
		if err != nil {
			telemetry.Log().Warn("skipping malformed key", zap.String("key_id", k.ID.String()), zap.Error(err))
			continue
		}
		jwks = append(jwks, jwk)
	}

	return map[string]any{"keys": jwks}, nil
}
