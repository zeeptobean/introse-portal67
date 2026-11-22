// Package dbx contains small PostgreSQL helpers shared by the feature stores.
package dbx

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
)

// WithTx runs fn in a transaction, committing when fn returns nil.
func WithTx(ctx context.Context, db *pgxpool.Pool, fn func(pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		if traceID := events.TraceID(ctx); traceID != "" {
			if _, err := tx.Exec(ctx, "SELECT set_config('portal67.trace_id', $1, true)", traceID); err != nil {
				return err
			}
		}
		return fn(tx)
	})
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// ContainsPattern returns a LIKE/ILIKE pattern that matches values containing
// query literally. Backslash, percent, and underscore are escaped with the
// default LIKE escape character, so they lose their wildcard meaning.
func ContainsPattern(query string) string {
	return "%" + likeEscaper.Replace(strings.TrimSpace(query)) + "%"
}
