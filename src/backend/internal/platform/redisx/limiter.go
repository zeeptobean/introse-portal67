package redisx

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type LoginLimiter struct {
	Client *redis.Client
	Prefix string
}

var checkLoginWindow = redis.NewScript(`
local n = redis.call('INCR', KEYS[1])
if n == 1 then redis.call('PEXPIRE', KEYS[1], ARGV[1]) end
return {n, redis.call('PTTL', KEYS[1])}
`)

func (l LoginLimiter) key(bucket, identifier string) string {
	prefix := l.Prefix
	if prefix == "" {
		prefix = "p67"
	}
	return fmt.Sprintf("%s:rl:login:%s:%x", prefix, bucket, sha256.Sum256([]byte(identifier)))
}

func (l LoginLimiter) Allow(ctx context.Context, bucket, identifier string, limit int, window time.Duration) (bool, time.Duration, error) {
	call, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	values, err := checkLoginWindow.Run(call, l.Client, []string{l.key(bucket, identifier)}, window.Milliseconds()).Int64Slice()
	if err != nil {
		return false, 0, err
	}
	return values[0] <= int64(limit), max(time.Millisecond*time.Duration(values[1]), time.Second), nil
}

func (l LoginLimiter) Reset(ctx context.Context, bucket, identifier string) {
	call, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	_ = l.Client.Del(call, l.key(bucket, identifier)).Err()
}
