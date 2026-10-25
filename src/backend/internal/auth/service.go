package auth

import (
	"context"
	"errors"
	"fmt"
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
	store *Store
	dummyHash string
}

func NewService(store *Store) *Service {
	dummyHash, err := argon2id.CreateHash("portal67-dummy-password", argon2id.DefaultParams)
	if err != nil {
		panic(fmt.Sprintf("create dummy password hash: %v", err))
	}
	return &Service{store: store, dummyHash: dummyHash}
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	c, err := s.store.credentialsByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		_, _ = argon2id.ComparePasswordAndHash(password, s.dummyHash)
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("authenticate: %w", err)
	}
	match, err := argon2id.ComparePasswordAndHash(password, c.PasswordHash)
	if err != nil {
		return User{}, fmt.Errorf("authenticate user %d: compare password: %w", c.ID, err)
	}
	if !match {
		return User{}, ErrInvalidCredentials
	}
	if !c.IsActive {
		return User{}, ErrInactive
	}
	return c.User, nil
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
	return s.store.listStudents(ctx, strings.TrimSpace(query))
}
