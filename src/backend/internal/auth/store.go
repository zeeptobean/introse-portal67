package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const userColumns = "id, email, institutional_id, full_name, date_of_birth, role, is_active, created_at"

// Store holds the database dependency for account and credential queries.
type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// credentials is a user together with their password hash. It never leaves
// the auth package so the hash cannot be serialized by accident.
type credentials struct {
	User
	PasswordHash string
}

// userByID returns pgx.ErrNoRows (wrapped) when no user has the given ID.
func (s *Store) userByID(ctx context.Context, id int64) (User, error) {
	rows, err := s.db.Query(ctx, "SELECT "+userColumns+" FROM users WHERE id = $1", id)
	if err != nil {
		return User{}, fmt.Errorf("find user %d: %w", id, err)
	}
	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return User{}, fmt.Errorf("find user %d: %w", id, err)
	}
	return u, nil
}

// credentialsByEmail returns pgx.ErrNoRows (wrapped) when no user has the
// given email.
func (s *Store) credentialsByEmail(ctx context.Context, email string) (credentials, error) {
	rows, err := s.db.Query(ctx, "SELECT "+userColumns+", password_hash FROM users WHERE email = $1", email)
	if err != nil {
		return credentials{}, fmt.Errorf("find credentials by email: %w", err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[credentials])
	if err != nil {
		return credentials{}, fmt.Errorf("find credentials by email: %w", err)
	}
	return c, nil
}

func (s *Store) updateUser(ctx context.Context, id int64, input UserUpdate) (User, error) {
	if input.DateOfBirth == nil && input.FullName == nil && input.InstitutionalID == nil {
		return User{}, errors.New("update user: no fields to update")
	}
	rows, err := s.db.Query(ctx, `
		UPDATE users
		SET full_name = COALESCE($2, full_name),
			institutional_id = COALESCE($3, institutional_id),
			date_of_birth = COALESCE($4, date_of_birth)
		WHERE id = $1
		RETURNING `+userColumns,
		id, input.FullName, input.InstitutionalID, input.DateOfBirth)
	if err != nil {
		return User{}, fmt.Errorf("update user %d: %w", id, err)
	}
	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return User{}, fmt.Errorf("update user %d: %w", id, err)
	}
	return u, nil
}
