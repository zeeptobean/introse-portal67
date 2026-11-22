package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/app"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/config"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events/redisbroker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		slog.Error("API stopped", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("configure PostgreSQL: %w", err)
	}
	defer pool.Close()

	var pingRedis func(context.Context) error
	if cfg.RedisURL != "" {
		broker, err := redisbroker.Open(cfg.RedisURL, cfg.RedisKeyPrefix)
		if err != nil {
			return err
		}
		defer broker.Close()
		pingRedis = broker.Ping
		relayCtx, stopRelay := context.WithCancel(ctx)
		relayDone := make(chan struct{})
		go func() {
			defer close(relayDone)
			events.NewRelay(pool, broker).Run(relayCtx)
		}()
		defer func() {
			stopRelay()
			<-relayDone
		}()
	}

	server := &http.Server{
		Addr:              cfg.APIAddr,
		Handler:           app.NewHandler(pool, cfg.TrustedProxyCIDRs, pingRedis),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("API listening", "address", cfg.APIAddr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown API: %w", err)
		}
		return nil
	}
}
