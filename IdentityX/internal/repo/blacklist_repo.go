package repo

import (
	"GoAuth/internal/models"
	"GoAuth/internal/sqlc"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type RevokedRefreshTokensRepo interface {
	Revoke(ctx context.Context, blacklist models.RefreshBlacklist) error
	RevokeMany(ctx context.Context, tokenIDs []uuid.UUID, expiresAts []time.Time) error
	GetByID(ctx context.Context, BlacklistID uuid.UUID) (*models.RefreshBlacklist, error)
	Delete(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type blacklistRepo struct {
	q   *sqlc.Queries
	log *zap.Logger
}

func NewRevokedRefreshTokensRepo(q *sqlc.Queries, l *zap.Logger) RevokedRefreshTokensRepo {
	return &blacklistRepo{
		q:   q,
		log: l,
	}
}

func copyBlacklistFromDB(dst *models.RefreshBlacklist, src *sqlc.RefreshBlacklist) error {
	return copier.Copy(dst, src)
}

func (b blacklistRepo) Revoke(ctx context.Context, blacklist models.RefreshBlacklist) error {
	if blacklist.TokenID == uuid.Nil {
		return errors.New("TokenID is not valid")
	}
	err := b.q.RevokeToken(ctx, sqlc.RevokeTokenParams{
		TokenID:   blacklist.TokenID,
		ExpiresAt: blacklist.ExpiresAt,
	})

	if err != nil {
		return err
	}

	return nil
}

func (b blacklistRepo) RevokeMany(ctx context.Context, tokenIDs []uuid.UUID, expiresAts []time.Time) error {
	_, err := b.q.RevokeManyTokens(ctx, sqlc.RevokeManyTokensParams{
		Column1: tokenIDs,
		Column2: expiresAts,
	})

	if err != nil {
		return err
	}

	return nil
}

func (b blacklistRepo) GetByID(ctx context.Context, BlacklistID uuid.UUID) (*models.RefreshBlacklist, error) {
	if BlacklistID == uuid.Nil {
		return nil, errors.New("BlacklistID is not valid")
	}

	sqlcRevokedToken, err := b.q.GetRevokedRefreshByID(ctx, BlacklistID)

	if err != nil {
		return nil, err
	}

	var revokedToken models.RefreshBlacklist
	if err = copyBlacklistFromDB(&revokedToken, &sqlcRevokedToken); err != nil {
		b.log.Error(
			"failed to copy revoked token",
			zap.Error(err),
			zap.String("session_id", sqlcRevokedToken.TokenID.String()),
		)
		return nil, fmt.Errorf("failed to copy revoked token: %w", err)
	}

	return &revokedToken, nil
}

func (b blacklistRepo) Delete(ctx context.Context, tokenID uuid.UUID) error {
	if tokenID == uuid.Nil {
		return errors.New("TokenID is not valid")
	}

	err := b.q.DeleteRevokedRefreshByID(ctx, tokenID)
	if err != nil {
		return err
	}

	return nil
}

func (b blacklistRepo) DeleteExpired(ctx context.Context) error {
	return b.q.DeleteExpiredRefreshTokens(ctx)
}
