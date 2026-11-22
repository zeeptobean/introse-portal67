package events

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Message is the wire payload of a committed outbox row.
type Message struct {
	Stream  string
	Type    string
	Payload string
}

// Publisher accepts a batch in order. If publication partially succeeds, the
// whole batch is eligible for retry using the original event IDs.
type Publisher interface {
	Publish(context.Context, []Message) error
}

type Relay struct {
	pool      *pgxpool.Pool
	publisher Publisher
}

func NewRelay(pool *pgxpool.Pool, publisher Publisher) *Relay {
	return &Relay{pool: pool, publisher: publisher}
}

func (r *Relay) Run(ctx context.Context) {
	for ctx.Err() == nil {
		attempt, cancel := context.WithTimeout(ctx, 4*time.Second)
		_, err := r.Flush(attempt)
		cancel()
		delay := 500 * time.Millisecond
		if err != nil && ctx.Err() == nil {
			slog.Warn("outbox relay will retry", "error", err)
			delay = 5 * time.Second
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
		}
	}
}

// Flush sends one batch. The transaction-level lock prevents overlapping
// workers from publishing newer messages before the current batch.
func (r *Relay) Flush(ctx context.Context) (int, error) {
	type row struct {
		Seq int64
		Message
	}
	var published int
	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		var locked bool
		if err := tx.QueryRow(ctx, "SELECT pg_try_advisory_xact_lock(hashtextextended(current_schema()::text, 6722))").Scan(&locked); err != nil || !locked {
			return err
		}
		rows, err := tx.Query(ctx, `
			SELECT seq, stream, type, payload::text
			FROM outbox WHERE published_at IS NULL
			ORDER BY seq LIMIT 100`)
		if err != nil {
			return err
		}
		batch, err := pgx.CollectRows(rows, func(record pgx.CollectableRow) (row, error) {
			var entry row
			err := record.Scan(&entry.Seq, &entry.Stream, &entry.Type, &entry.Payload)
			return entry, err
		})
		if err != nil || len(batch) == 0 {
			return err
		}
		messages := make([]Message, 0, len(batch))
		sequences := make([]int64, 0, len(batch))
		for _, entry := range batch {
			messages = append(messages, entry.Message)
			sequences = append(sequences, entry.Seq)
		}
		if err := r.publisher.Publish(ctx, messages); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, "UPDATE outbox SET published_at = clock_timestamp() WHERE seq = ANY($1)", sequences)
		published = int(tag.RowsAffected())
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("flush event outbox: %w", err)
	}
	return published, nil
}

// Replay queues stored events again with their original event IDs.
func Replay(ctx context.Context, pool *pgxpool.Pool, since time.Time) (int64, error) {
	var count int64
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock(hashtextextended(current_schema()::text, 6722))"); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, "UPDATE outbox SET published_at = NULL WHERE published_at IS NOT NULL AND created_at >= $1", since)
		count = tag.RowsAffected()
		return err
	})
	return count, err
}
