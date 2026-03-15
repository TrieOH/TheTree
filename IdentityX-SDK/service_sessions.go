package goauth

import (
	"context"
	"time"

	"github.com/MintzyG/fail/v3"
)

type SessionCache interface {
	GetSession(ctx context.Context, id string) ([]byte, error)
	SetSession(ctx context.Context, id string, data []byte, ttl time.Duration) error
	DeleteSession(ctx context.Context, id string) error
}

type SessionRuntime struct {
	cache SessionCache
}

func (r *SessionRuntime) Create(
	ctx context.Context,
	sessionID string,
	payload []byte,
	ttl time.Duration,
) error {
	return r.cache.SetSession(ctx, sessionID, payload, ttl)
}

func (r *SessionRuntime) Get(
	ctx context.Context,
	sessionID string,
) ([]byte, error) {
	return r.cache.GetSession(ctx, sessionID)
}

func (r *SessionRuntime) Delete(
	ctx context.Context,
	sessionID string,
) error {
	return r.cache.DeleteSession(ctx, sessionID)
}

type SnapshotBuilder func(claims *AccessClaims) ([]byte, error)

type SessionResult struct {
	SessionID string        `json:"session_id"`
	TTL       time.Time     `json:"ttl"`
	Claims    *AccessClaims `json:"claims"`
}

func (c *Client) ExchangeAndCreateSession(
	ctx context.Context,
	accessToken string,
	build SnapshotBuilder,
) (*SessionResult, error) {
	if c.Sessions == nil {
		return nil, fail.New(SDKSessionsNotConfiguredID)
	}

	claims, err := c.Tokens.VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	snapshot, err := build(claims)
	if err != nil {
		return nil, fail.New(SDKSnapshotBuildFailedID).
			WithArgs(err.Error())
	}

	ttl := time.Until(claims.ExpiresAt.Time)

	err = c.Sessions.Create(
		ctx,
		claims.Sub.SessionID.String(),
		snapshot,
		ttl,
	)
	if err != nil {
		return nil, err
	}

	return &SessionResult{
		SessionID: claims.Sub.SessionID.String(),
		TTL:       claims.ExpiresAt.Time,
		Claims:    claims,
	}, nil
}
