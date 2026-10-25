package auth

import (
	"context"
	"fmt"

	"github.com/alexedwards/argon2id"
)

// hashParams follows the OWASP argon2id baseline (19 MiB, 2 passes, 1 lane).
// argon2id.DefaultParams uses 64 MiB per hash, which lets every login attempt
// claim 64 MiB of memory.
var hashParams = &argon2id.Params{
	Memory:      19 * 1024,
	Iterations:  2,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

// maxConcurrentHashes bounds the memory held by password hashing at once,
// whatever the request rate: at most this many hashes run at the same time.
const maxConcurrentHashes = 4

// hasher runs password hashes with bounded concurrency.
type hasher struct {
	slots chan struct{}
}

func newHasher() *hasher {
	return &hasher{slots: make(chan struct{}, maxConcurrentHashes)}
}

// acquire waits for a free hashing slot. The caller must call release.
func (h *hasher) acquire(ctx context.Context) error {
	select {
	case h.slots <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for password hashing slot: %w", ctx.Err())
	}
}

func (h *hasher) release() { <-h.slots }

func (h *hasher) compare(ctx context.Context, password, hash string) (bool, error) {
	if err := h.acquire(ctx); err != nil {
		return false, err
	}
	defer h.release()
	return argon2id.ComparePasswordAndHash(password, hash)
}

func (h *hasher) create(ctx context.Context, password string) (string, error) {
	if err := h.acquire(ctx); err != nil {
		return "", err
	}
	defer h.release()
	return argon2id.CreateHash(password, hashParams)
}

// needsRehash reports whether hash was created with parameters other than
// hashParams, such as the 64 MiB defaults used by older accounts.
func needsRehash(hash string) bool {
	params, _, _, err := argon2id.DecodeHash(hash)
	if err != nil {
		return false
	}
	return params.Memory != hashParams.Memory ||
		params.Iterations != hashParams.Iterations ||
		params.Parallelism != hashParams.Parallelism ||
		params.SaltLength != hashParams.SaltLength ||
		params.KeyLength != hashParams.KeyLength
}
