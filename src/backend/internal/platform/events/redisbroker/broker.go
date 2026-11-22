package redisbroker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
)

type Broker struct {
	client *redis.Client
	prefix string
}

func Open(rawURL, prefix string) (*Broker, error) {
	options, err := redis.ParseURL(rawURL)
	if err != nil || options.Network != "tcp" {
		// Parsing errors can include the URL, including its credentials.
		return nil, errors.New("REDIS_URL must be a valid redis:// or rediss:// URL")
	}
	if prefix == "" {
		prefix = "p67"
	}
	options.DialTimeout = 2 * time.Second
	options.ReadTimeout = 2 * time.Second
	options.WriteTimeout = 2 * time.Second
	options.ContextTimeoutEnabled = true
	options.MaxRetries = -1 // The outbox owns retries.
	return &Broker{client: redis.NewClient(options), prefix: prefix}, nil
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
				Stream: b.prefix + ":ev:" + message.Stream,
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
