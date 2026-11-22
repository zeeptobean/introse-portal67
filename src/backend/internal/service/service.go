package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/app"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/config"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/notification"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/authn"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events/redisbroker"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/redisx"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/server"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/projections"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
)

func Main(kind, defaultPort string) {
	server.Main(kind, func(ctx context.Context) error { return run(ctx, kind, defaultPort) })
}

func OpenDB(ctx context.Context, databaseURL, kind string) (*pgxpool.Pool, error) {
	if _, ok := routing.Routes[kind]; !ok {
		return nil, errors.New("unknown service")
	}
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, errors.New("invalid service DATABASE_URL")
	}
	role := "portal67_" + kind
	if cfg.ConnConfig.User != role {
		return nil, fmt.Errorf("DATABASE_URL must use the %s role", role)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = kind
	// Do not allow a DSN option to replace the service's configured search path.
	delete(cfg.ConnConfig.RuntimeParams, "options")
	return pgxpool.NewWithConfig(ctx, cfg)
}

func run(ctx context.Context, kind, defaultPort string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cfg.RedisURL == "" {
		return errors.New("REDIS_URL is required for split services")
	}
	keys, err := authn.ParseVerificationKeys(os.Getenv("IDENTITY_PUBLIC_KEYS"))
	if err != nil {
		return err
	}
	pool, err := OpenDB(ctx, cfg.DatabaseURL, kind)
	if err != nil {
		return err
	}
	defer pool.Close()
	broker, err := redisbroker.Open(cfg.RedisURL, cfg.RedisKeyPrefix)
	if err != nil {
		return err
	}
	defer broker.Close()
	checks := []func(context.Context) error{broker.Ping, queueReady(pool)}
	var handler http.Handler
	if kind == routing.Identity {
		client, err := redisx.Open(cfg.RedisURL)
		if err != nil {
			return err
		}
		defer client.Close()
		limiter := redisx.LoginLimiter{Client: client, Prefix: cfg.RedisKeyPrefix}
		handler = app.NewIdentityHandler(pool, cfg.TrustedProxyCIDRs, keys, limiter, checks...)
	} else {
		handler = app.NewDomainHandler(kind, pool, keys, checks...)
	}

	workerCtx, stop := context.WithCancel(ctx)
	var workers sync.WaitGroup
	defer func() {
		stop()
		workers.Wait()
	}()
	if kind != routing.Notification {
		workers.Go(func() { events.NewRelay(pool, broker).Run(workerCtx) })
	}
	host, _ := os.Hostname()
	worker := fmt.Sprintf("%s-%s-%d", kind, host, os.Getpid())
	consumes := projections.For(kind)
	if kind == routing.Notification {
		consumes = map[string]events.Handler{
			"identity": projections.User, "academic": notification.Project,
			"affairs": notification.Project, "registration": notification.Project,
		}
	}
	for stream, project := range consumes {
		consumer := events.Consumer{Broker: broker, Group: kind,
			Stream: stream, Worker: worker, Apply: events.Transactional(pool, stream, project)}
		workers.Go(func() { consumer.Run(workerCtx) })
	}

	addr := os.Getenv("SERVICE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:" + defaultPort
	}
	return server.Listen(ctx, kind, addr, handler)
}

func queueReady(pool *pgxpool.Pool) func(context.Context) error {
	return func(ctx context.Context) error {
		var behind bool
		err := pool.QueryRow(ctx, `
			SELECT EXISTS (SELECT 1 FROM outbox WHERE published_at IS NULL
				AND created_at < clock_timestamp() - interval '60 seconds')`).Scan(&behind)
		if err != nil {
			return err
		}
		if behind {
			return errors.New("event outbox is behind")
		}
		return nil
	}
}
