package academics

import "github.com/jackc/pgx/v5/pgxpool"

// Store holds the database dependency for course registration, scheduling, and academic results.
type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// TODO: Add parameterized SQL queries and row scanning here.
// Transactional operations must share the pgx.Tx owned by the service.
