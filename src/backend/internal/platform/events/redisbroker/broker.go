package redisbroker

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/redisx"
)

type Broker struct {
	client *redis.Client
	prefix string
}

func Open(rawURL, prefix string) (*Broker, error) {
	client, err := redisx.Open(rawURL)
	if err != nil {
		return nil, err
	}
	if prefix == "" {
		prefix = "p67"
	}
	return &Broker{client: client, prefix: prefix}, nil
}

func (b *Broker) Close() error {
	return b.client.Close()
}

func (b *Broker) Ping(ctx context.Context) error {
	return b.client.Ping(ctx).Err()
}

func (b *Broker) Publish(ctx context.Context, batch []events.Message) error {
	_, err := b.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, message := range batch {
			pipe.XAdd(ctx, &redis.XAddArgs{
				Stream: b.streamKey(message.Stream),
				Values: []string{"type", message.Type, "payload", message.Payload},
			})
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("publish Redis stream batch: %w", err)
	}
	return nil
}
