// Package events stores domain events in PostgreSQL and relays them to the broker.
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Aggregate struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Version int64  `json:"version"`
}

type Event struct {
	ID            uuid.UUID       `json:"id"`
	Type          string          `json:"type"`
	SchemaVersion int             `json:"schema_version"`
	OccurredAt    time.Time       `json:"occurred_at"`
	Producer      string          `json:"producer"`
	Aggregate     Aggregate       `json:"aggregate"`
	TraceID       string          `json:"trace_id"`
	ActorID       int64           `json:"actor_id"`
	Data          json.RawMessage `json:"data"`
}

type traceIDKey struct{}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" || len(traceID) > 128 {
		return ctx
	}
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func TraceID(ctx context.Context) string {
	traceID, _ := ctx.Value(traceIDKey{}).(string)
	return traceID
}

// Add saves the event in the caller's business transaction.
func Add(ctx context.Context, tx pgx.Tx, stream string, event Event) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate event ID: %w", err)
	}
	event.ID = id
	event.TraceID = TraceID(ctx)
	if event.TraceID == "" {
		event.TraceID = id.String()
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("encode %s event: %w", event.Type, err)
	}
	_, err = tx.Exec(ctx, "INSERT INTO outbox (event_id, stream, type, payload) VALUES ($1, $2, $3, $4)",
		event.ID.String(), stream, event.Type, payload)
	if err != nil {
		return fmt.Errorf("append %s event: %w", event.Type, err)
	}
	return nil
}
