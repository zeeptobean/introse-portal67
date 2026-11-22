package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/config"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/gateway"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/authn"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events/redisbroker"
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
		prefix := os.Getenv("REDIS_KEY_PREFIX")
		broker, err := redisbroker.Open(raw, prefix)
		if err != nil {
			return err
		}
		defer broker.Close()
		workerCtx, cancel := context.WithCancel(ctx)
		workerDone := make(chan struct{})
		consumer := events.Consumer{Broker: broker, Stream: "identity", Group: "gateway",
			Worker: fmt.Sprintf("gateway-%d", os.Getpid()), Apply: gateway.InvalidateCachedUsers(cache, prefix)}
		go func() {
			defer close(workerDone)
			consumer.Run(workerCtx)
		}()
		defer func() {
			cancel()
			<-workerDone
		}()
	}
	endpoints := make(map[string]string, len(routing.Routes))
	for service := range routing.Routes {
		endpoints[service] = os.Getenv(strings.ToUpper(service) + "_URL")
	}
	handler, err := gateway.New(key, endpoints, trusted, cache, os.Getenv("REDIS_KEY_PREFIX"))
	if err != nil {
		return err
	}
	addr := os.Getenv("SERVICE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	return server.Listen(ctx, "gateway", addr, handler)
}
