package purchases

import (
	"context"
	"encoding/json"
	"time"
	"univents/internal/shared/contracts"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CheckoutSessionStore struct {
	rdb *redis.Client
}

func NewCheckoutSessionStore(rdb *redis.Client) *CheckoutSessionStore {
	return &CheckoutSessionStore{rdb: rdb}
}

func (s *CheckoutSessionStore) Save(ctx context.Context, session contracts.CheckoutSession) error {
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	ttl := time.Until(session.ExpiresAt) + 2*time.Minute
	return s.rdb.Set(ctx, contracts.CheckoutSessionKey(session.UserID, session.SessionID), b, ttl).Err()
}

func (s *CheckoutSessionStore) Load(ctx context.Context, userID, sessionID uuid.UUID) (*contracts.CheckoutSession, error) {
	b, err := s.rdb.Get(ctx, contracts.CheckoutSessionKey(userID, sessionID)).Bytes()
	if err != nil {
		return nil, err
	}
	var session contracts.CheckoutSession
	if err = json.Unmarshal(b, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *CheckoutSessionStore) Delete(ctx context.Context, userID, sessionID uuid.UUID) error {
	return s.rdb.Del(ctx, contracts.CheckoutSessionKey(userID, sessionID)).Err()
}
