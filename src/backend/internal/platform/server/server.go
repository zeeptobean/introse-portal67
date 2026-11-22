package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
)

func Main(service string, run func(context.Context) error) {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("load environment", "service", service, "error", err)
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx); err != nil {
		slog.Error("service stopped", "service", service, "error", err)
		os.Exit(1)
	}
}

func HealthRoutes(mux *http.ServeMux, checks ...func(context.Context) error) {
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		for _, check := range checks {
			if err := check(ctx); err != nil {
				httpx.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
				return
			}
		}
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
}

func Listen(ctx context.Context, service, addr string, handler http.Handler) error {
	srv := &http.Server{
		Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: time.Minute,
	}
	stopped := make(chan error, 1)
	go func() {
		slog.Info("listening", "service", service, "address", addr)
		stopped <- srv.ListenAndServe()
	}()
	select {
	case err := <-stopped:
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			_ = srv.Close()
			return err
		}
		return nil
	}
}
