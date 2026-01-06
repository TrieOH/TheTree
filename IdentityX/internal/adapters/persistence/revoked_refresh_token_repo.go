package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/revoked_refreshes"
	"GoAuth/internal/ports/outbound"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type revokedRefreshTokensRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *revokedRefreshTokensRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(txKeyValue).(*sql.Tx); ok {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbound.RevokedRefreshTokenRepository = (*revokedRefreshTokensRepo)(nil)

func NewRevokedRefreshTokensRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.RevokedRefreshTokenRepository {
	return &revokedRefreshTokensRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func mapRevokedRefreshTokenFromDB(dst *revoked_refreshes.RevokedRefreshToken, src *sqlc.RevokedRefreshToken) {
	dst.TokenID = src.TokenID
	dst.CreatedAt = src.CreatedAt
	dst.ExpiresAt = src.ExpiresAt
}

func (repo *revokedRefreshTokensRepo) Revoke(ctx context.Context, toRevoke revoked_refreshes.RevokedRefreshToken) error {
	ctx, span := repo.tracer.Start(ctx, "RevokedRefreshTokensRepo.Revoke",
		trace.WithAttributes(
			attribute.String("revoked_token.id", toRevoke.TokenID.String()),
			attribute.Int64("revoked_token.expires_at", toRevoke.ExpiresAt.Unix()),
		),
	)
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("revoke.success", err == nil))
		}
	}()

	err = repo.queries(ctx).RevokeToken(ctx, sqlc.RevokeTokenParams{
		TokenID:   toRevoke.TokenID,
		ExpiresAt: toRevoke.ExpiresAt,
	})

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return sqlErr
	}

	return nil
}

func (repo *revokedRefreshTokensRepo) RevokeMany(ctx context.Context, tokenIDs []uuid.UUID, expiresAts []time.Time) error {
	ctx, span := repo.tracer.Start(ctx, "RevokedRefreshTokensRepo.RevokeMany",
		trace.WithAttributes(
			attribute.Int("revoke_many.count", len(tokenIDs)),
		),
	)
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("revoke_many.success", err == nil))
		}
	}()

	var revokedTokens []uuid.UUID
	revokedTokens, err = repo.queries(ctx).RevokeManyTokens(ctx, sqlc.RevokeManyTokensParams{
		Column1: tokenIDs,
		Column2: expiresAts,
	})
	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return sqlErr
	}

	span.SetAttributes(attribute.Int("revoke_many.revoked.count", len(revokedTokens)))

	return nil
}

func (repo *revokedRefreshTokensRepo) GetByID(ctx context.Context, revokedID uuid.UUID) (*revoked_refreshes.RevokedRefreshToken, error) {
	ctx, span := repo.tracer.Start(ctx, "RevokedRefreshTokenRepo.GetByID",
		trace.WithAttributes(
			attribute.String("revoked_token_id", revokedID.String()),
		),
	)
	defer span.End()

	sqlcRevokedToken, err := repo.queries(ctx).GetRevokedRefreshByID(ctx, revokedID)

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
	}

	span.SetAttributes(
		attribute.Int64("revoked_at", sqlcRevokedToken.CreatedAt.Unix()),
		attribute.Int64("expires_at", sqlcRevokedToken.ExpiresAt.Unix()),
	)

	var revokedToken revoked_refreshes.RevokedRefreshToken
	mapRevokedRefreshTokenFromDB(&revokedToken, &sqlcRevokedToken)

	return &revokedToken, nil
}

func (repo *revokedRefreshTokensRepo) Delete(ctx context.Context, tokenID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RevokedRefreshTokenRepo.Delete",
		trace.WithAttributes(
			attribute.String("deleted_revoked_token_id", tokenID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).DeleteRevokedRefreshByID(ctx, tokenID)
	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return sqlErr
	}

	return nil
}

func (repo *revokedRefreshTokensRepo) DeleteExpired(ctx context.Context) error {
	ctx, span := repo.tracer.Start(ctx, "RevokedRefreshTokensRepo.DeleteExpired")
	defer span.End()

	if err := repo.queries(ctx).DeleteExpiredRefreshTokens(ctx); err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return sqlErr
	}

	return nil
}

func (repo *revokedRefreshTokensRepo) IsRevoked(ctx context.Context, tokenID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "RevokedRefreshTokenRepo.IsRevoked",
		trace.WithAttributes(
			attribute.String("revoked_token.id", tokenID.String()),
		),
	)
	defer span.End()

	isRevoked, err := repo.queries(ctx).IsRefreshTokenRevoked(ctx, tokenID)

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return false, sqlErr
	}

	span.SetAttributes(attribute.Bool("is_revoked", isRevoked))
	return isRevoked, nil
}
