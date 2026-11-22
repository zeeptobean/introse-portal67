package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Delivery struct {
	Message
	ID       string
	Attempts int64
}

type Subscriber interface {
	CreateGroup(ctx context.Context, stream, group string) error
	Read(ctx context.Context, stream, group, worker string) (Delivery, error)
	Ack(ctx context.Context, delivery Delivery, group string) error
	DeadLetter(ctx context.Context, delivery Delivery, group string, reason string) error
}

type Handler func(context.Context, pgx.Tx, Event) error

type Consumer struct {
	Broker Subscriber
	Stream string
	Group  string
	Worker string
	Apply  func(context.Context, Event) error
}

func Transactional(pool *pgxpool.Pool, stream string, handle Handler) func(context.Context, Event) error {
	return func(ctx context.Context, event Event) error {
		return pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
			tag, err := tx.Exec(ctx, `
				INSERT INTO inbox (event_id, stream, type) VALUES ($1, $2, $3)
				ON CONFLICT (event_id) DO NOTHING`, event.ID.String(), stream, event.Type)
			if err != nil || tag.RowsAffected() == 0 {
				return err
			}
			return handle(ctx, tx, event)
		})
	}
}

func (c Consumer) Run(ctx context.Context) {
	for ctx.Err() == nil {
		if err := c.Broker.CreateGroup(ctx, c.Stream, c.Group); err != nil {
			c.waitAfterError(ctx, err)
			continue
		}
		for ctx.Err() == nil {
			if err := c.Step(ctx); err != nil {
				c.waitAfterError(ctx, err)
				break
			}
		}
	}
}

func (c Consumer) waitAfterError(ctx context.Context, err error) {
	if ctx.Err() != nil {
		return
	}
	slog.Warn("stream consumer will retry", "stream", c.Stream, "group", c.Group, "error", err)
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
	}
}

func (c Consumer) Step(ctx context.Context) error {
	delivery, err := c.Broker.Read(ctx, c.Stream, c.Group, c.Worker)
	if err != nil || delivery.ID == "" {
		return err
	}
	if err := c.process(ctx, delivery); err != nil {
		if delivery.Attempts < 5 || ctx.Err() != nil {
			return fmt.Errorf("process %s: %w", delivery.ID, err)
		}
		if err := c.Broker.DeadLetter(ctx, delivery, c.Group, err.Error()); err != nil {
			return fmt.Errorf("dead-letter %s: %w", delivery.ID, err)
		}
		slog.Error("event sent to dead letter stream", "stream", c.Stream, "group", c.Group, "entry_id", delivery.ID)
		return nil
	}
	return c.Broker.Ack(ctx, delivery, c.Group)
}

func (c Consumer) process(ctx context.Context, delivery Delivery) error {
	var event Event
	if err := json.Unmarshal([]byte(delivery.Payload), &event); err != nil {
		return fmt.Errorf("invalid event JSON: %w", err)
	}
	if event.ID == uuid.Nil || event.Type == "" || event.Type != delivery.Type || event.SchemaVersion < 1 {
		return fmt.Errorf("invalid event envelope")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return c.Apply(ctx, event)
}
