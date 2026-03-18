package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"univents/internal/commerce/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type redisInventoryPublisher struct {
	redis *redis.Client
}

func NewRedisInventoryPublisher(redis *redis.Client) domain.InventoryPublisher {
	return &redisInventoryPublisher{redis: redis}
}

func inventoryChannel(editionID uuid.UUID) string {
	return fmt.Sprintf("inventory:%s", editionID)
}

func (p *redisInventoryPublisher) Publish(ctx context.Context, editionID uuid.UUID, updates []domain.InventoryUpdate) error {
	payload, err := json.Marshal(updates)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory updates: %w", err)
	}
	return p.redis.Publish(ctx, inventoryChannel(editionID), payload).Err()
}

type redisInventorySubscriber struct {
	redis *redis.Client
}

func NewRedisInventorySubscriber(redis *redis.Client) domain.InventorySubscriber {
	return &redisInventorySubscriber{redis: redis}
}

func (s *redisInventorySubscriber) Subscribe(ctx context.Context, editionID uuid.UUID) (<-chan []domain.InventoryUpdate, error) {
	ch := make(chan []domain.InventoryUpdate, 16)
	sub := s.redis.Subscribe(ctx, inventoryChannel(editionID))

	// verify subscription succeeded
	if _, err := sub.Receive(ctx); err != nil {
		_ = sub.Close()
		return nil, fmt.Errorf("failed to subscribe to inventory channel: %w", err)
	}

	go func() {
		defer close(ch)
		defer sub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-sub.Channel():
				if !ok {
					return
				}
				var updates []domain.InventoryUpdate
				if err := json.Unmarshal([]byte(msg.Payload), &updates); err != nil {
					continue
				}
				select {
				case ch <- updates:
				default:
					// client too slow, drop the message
				}
			}
		}
	}()

	return ch, nil
}
