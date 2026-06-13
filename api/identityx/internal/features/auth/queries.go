package auth

import (
	"context"

	"IdentityX/contracts"
	"IdentityX/internal/shared/ports"
	"lib/database"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueryService struct {
	keys   ports.KeysRepository
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewQueryService(
	Keys ports.KeysRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		keys:   Keys,
		logger: logger,
		tracer: tracer,
		tx:     tx,
	}
}

func (uc *QueryService) GetJWKS(ctx context.Context, projectID *uuid.UUID) (map[string]any, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.GetJWKS")
	defer span.End()

	keys, err := uc.keys.ListPublicKeys(ctx, projectID)
	if err != nil {
		telemetry.Log().Error("Failed listing public token keys", zap.Error(err), zap.Any("project_id", projectID))
		return nil, fun.Err("JWKS retrieval failed").Internal()
	}

	jwkKeys := make([]any, 0, len(keys))
	for _, k := range keys {
		var jwk map[string]any
		jwk, err = contracts.PublicKeyToJWK(k)
		if err != nil {
			return nil, err
		}
		jwkKeys = append(jwkKeys, jwk)
	}

	return map[string]any{"keys": jwkKeys}, nil
}
