// Package dbx contains small PostgreSQL helpers shared by the feature stores.
package dbx

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithTx runs fn in a transaction, committing when fn returns nil.
func WithTx(ctx context.Context, db *pgxpool.Pool, fn func(pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, db, fn)
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// ContainsPattern returns a LIKE/ILIKE pattern that matches values containing
// query literally. Backslash, percent, and underscore are escaped with the
// default LIKE escape character, so they lose their wildcard meaning.
func ContainsPattern(query string) string {
	return "%" + likeEscaper.Replace(strings.TrimSpace(query)) + "%"
}
