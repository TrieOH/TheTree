package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/revoked_refreshes"
	"GoAuth/internal/ports/outbound"
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type revokedRefreshTokensRepo struct {
	q   *sqlc.Queries
	log *zap.Logger
}

var _ outbound.RevokedRefreshTokenRepository = (*revokedRefreshTokensRepo)(nil)

func NewRevokedRefreshTokensRepo(q *sqlc.Queries, l *zap.Logger) outbound.RevokedRefreshTokenRepository {
	return &revokedRefreshTokensRepo{
		q:   q,
		log: l,
	}
}

func mapRevokedRefreshTokenFromDB(dst *revoked_refreshes.RevokedRefreshToken, src *sqlc.RevokedRefreshToken) {
	dst.TokenID = src.TokenID
	dst.CreatedAt = src.CreatedAt
	dst.ExpiresAt = src.ExpiresAt
}

func (r revokedRefreshTokensRepo) Revoke(ctx context.Context, blacklist revoked_refreshes.RevokedRefreshToken) error {
	ctx, span := GoAuthRepoTracer.Start(ctx, "RevokedRefreshTokensRepo.Revoke",
		trace.WithAttributes(
			attribute.String("revoked_token.id", blacklist.TokenID.String()),
			attribute.Int64("revoked_token.expires_at", blacklist.ExpiresAt.Unix()),
		),
	)
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("revoke.success", err == nil))
		}
	}()

	err = r.q.RevokeToken(ctx, sqlc.RevokeTokenParams{
		TokenID:   blacklist.TokenID,
		ExpiresAt: blacklist.ExpiresAt,
	})

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return sqlErr
	}

	return nil
}

func (r revokedRefreshTokensRepo) RevokeMany(ctx context.Context, tokenIDs []uuid.UUID, expiresAts []time.Time) error {
	ctx, span := GoAuthRepoTracer.Start(ctx, "RevokedRefreshTokensRepo.RevokeMany",
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
	revokedTokens, err = r.q.RevokeManyTokens(ctx, sqlc.RevokeManyTokensParams{
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

func (r revokedRefreshTokensRepo) GetByID(ctx context.Context, RevokedTokenID uuid.UUID) (*revoked_refreshes.RevokedRefreshToken, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "RevokedRefreshTokenRepo.GetByID",
		trace.WithAttributes(
			attribute.String("revoked_token_id", RevokedTokenID.String()),
		),
	)
	defer span.End()

	sqlcRevokedToken, err := r.q.GetRevokedRefreshByID(ctx, RevokedTokenID)

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

func (r revokedRefreshTokensRepo) Delete(ctx context.Context, tokenID uuid.UUID) error {
	ctx, span := GoAuthRepoTracer.Start(ctx, "RevokedRefreshTokenRepo.Delete",
		trace.WithAttributes(
			attribute.String("deleted_revoked_token_id", tokenID.String()),
		),
	)
	defer span.End()

	err := r.q.DeleteRevokedRefreshByID(ctx, tokenID)
	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return sqlErr
	}

	return nil
}

func (r revokedRefreshTokensRepo) DeleteExpired(ctx context.Context) error {
	return r.q.DeleteExpiredRefreshTokens(ctx)
}

func (r revokedRefreshTokensRepo) IsRevoked(ctx context.Context, tokenID uuid.UUID) (bool, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "RevokedRefreshTokenRepo.IsRevoked",
		trace.WithAttributes(
			attribute.String("revoked_token.id", tokenID.String()),
		),
	)
	defer span.End()

	isRevoked, err := r.q.IsRefreshTokenRevoked(ctx, tokenID)

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return false, sqlErr
	}

	span.SetAttributes(attribute.Bool("is_revoked", isRevoked))
	return isRevoked, nil
}
