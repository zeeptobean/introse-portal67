package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInactive           = errors.New("account is inactive")
	ErrUserNotFound       = errors.New("user not found")
)

type Service struct {
	store     *Store
	hasher    *hasher
	dummyHash string
}

func NewService(store *Store) *Service {
	dummyHash, err := argon2id.CreateHash("portal67-dummy-password", hashParams)
	if err != nil {
		panic(fmt.Sprintf("create dummy password hash: %v", err))
	}
	return &Service{store: store, hasher: newHasher(), dummyHash: dummyHash}
}

// NormalizeEmail returns the form in which emails are stored.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// Authenticate checks the password. A correct password stored with outdated
// hash parameters is rehashed with the current ones.
func (s *Service) Authenticate(ctx context.Context, email, password string) (User, error) {
	c, err := s.store.credentialsByEmail(ctx, NormalizeEmail(email))
	if errors.Is(err, pgx.ErrNoRows) {
		// Spend the same work as a real check so response times do not
		// reveal which emails exist.
		if _, err := s.hasher.compare(ctx, password, s.dummyHash); err != nil {
			return User{}, fmt.Errorf("authenticate: %w", err)
		}
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("authenticate: %w", err)
	}
	match, err := s.hasher.compare(ctx, password, c.PasswordHash)
	if err != nil {
		return User{}, fmt.Errorf("authenticate user %d: compare password: %w", c.ID, err)
	}
	if !match {
		return User{}, ErrInvalidCredentials
	}
	if needsRehash(c.PasswordHash) {
		s.rehash(ctx, c.ID, password)
	}
	if !c.IsActive {
		return User{}, ErrInactive
	}
	return c.User, nil
}

// rehash stores the password with the current hash parameters. A failure is
// logged and retried on the next login.
func (s *Service) rehash(ctx context.Context, userID int64, password string) {
	hash, err := s.hasher.create(ctx, password)
	if err == nil {
		err = s.store.updatePasswordHash(ctx, userID, hash)
	}
	if err != nil {
		slog.WarnContext(ctx, "rehash password", "user_id", userID, "error", err)
	}
}

// UserByID returns ErrUserNotFound when no user has the given ID.
func (s *Service) UserByID(ctx context.Context, id int64) (User, error) {
	u, err := s.store.userByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	return u, err
}

// Students searches students by email, name, or institutional ID.
func (s *Service) Students(ctx context.Context, query string) ([]User, error) {
	return s.store.listStudents(ctx, query)
}
