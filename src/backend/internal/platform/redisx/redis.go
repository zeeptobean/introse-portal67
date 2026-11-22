package redisx

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func Open(rawURL string) (*redis.Client, error) {
	options, err := redis.ParseURL(rawURL)
	if err != nil || options.Network != "tcp" {
		// Parsing errors can include the URL, including its credentials.
		return nil, errors.New("REDIS_URL must be a valid redis:// or rediss:// URL")
	}
	options.DialTimeout = 2 * time.Second
	options.ReadTimeout = 2 * time.Second
	options.WriteTimeout = 2 * time.Second
	options.ContextTimeoutEnabled = true
	options.MaxRetries = -1
	return redis.NewClient(options), nil
}
