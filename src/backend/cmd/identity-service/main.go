package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/app"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/config"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/authn"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/redisx"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/server"
)

func main() { server.Main("identity", run) }

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	keys, err := authn.ParseVerificationKeys(os.Getenv("IDENTITY_PUBLIC_KEYS"))
	if err != nil {
		return err
	}
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	var limiter auth.SharedLoginLimiter
	if cfg.RedisURL != "" {
		client, err := redisx.Open(cfg.RedisURL)
		if err != nil {
			return err
		}
		defer client.Close()
		limiter = redisx.LoginLimiter{Client: client, Prefix: cfg.RedisKeyPrefix}
	}
	addr := os.Getenv("SERVICE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8081"
	}
	return server.Listen(ctx, "identity", addr, app.NewIdentityHandler(pool, cfg.TrustedProxyCIDRs, keys, limiter))
}
