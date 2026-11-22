package main

import (
	"context"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/config"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/gateway"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/authn"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/redisx"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/server"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
)

func main() { server.Main("gateway", run) }

func run(ctx context.Context) error {
	key, err := authn.ParseSigningKey(os.Getenv("IDENTITY_SIGNING_KEY"))
	if err != nil {
		return err
	}
	trusted, err := config.TrustedProxyCIDRs()
	if err != nil {
		return err
	}
	var cache *redis.Client
	if raw := os.Getenv("REDIS_URL"); raw != "" {
		cache, err = redisx.Open(raw)
		if err != nil {
			return err
		}
		defer cache.Close()
	}
	endpoints := make(map[string]string, len(routing.Routes))
	for service := range routing.Routes {
		endpoints[service] = os.Getenv(strings.ToUpper(service) + "_URL")
	}
	handler, err := gateway.New(key, endpoints, os.Getenv("LEGACY_API_URL"), trusted, cache, os.Getenv("REDIS_KEY_PREFIX"))
	if err != nil {
		return err
	}
	addr := os.Getenv("SERVICE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	return server.Listen(ctx, "gateway", addr, handler)
}
