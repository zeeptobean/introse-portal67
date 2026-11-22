package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events/redisbroker"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/server"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/service"
)

func main() { server.Main("portalctl", run) }

func run(ctx context.Context) error {
	args := os.Args[1:]
	if len(args) == 0 {
		return errors.New("usage: portalctl keygen | migrate | replay <RFC3339> | republish | dlq-retry <stream> <group> <entry-id>")
	}
	switch args[0] {
	case "keygen":
		public, private, err := ed25519.GenerateKey(nil)
		if err != nil {
			return err
		}
		fmt.Printf("IDENTITY_SIGNING_KEY=%s\nIDENTITY_PUBLIC_KEYS=%s\n",
			base64.RawURLEncoding.EncodeToString(private.Seed()), base64.RawURLEncoding.EncodeToString(public))
		return nil
	case "migrate":
		configured := map[string]string{}
		for name := range routing.Routes {
			configured[name] = os.Getenv(strings.ToUpper(name) + "_DATABASE_URL")
		}
		return service.Migrate(ctx, os.Getenv("PORTAL67_ADMIN_DATABASE_URL"), "db", os.Getenv("PORTAL67_SEED_DEMO") == "1", configured)
	case "dlq-retry":
		if len(args) != 4 {
			return errors.New("usage: portalctl dlq-retry <stream> <group> <entry-id>")
		}
		broker, err := redisbroker.Open(os.Getenv("REDIS_URL"), os.Getenv("REDIS_KEY_PREFIX"))
		if err != nil {
			return err
		}
		defer broker.Close()
		found, err := broker.RetryDeadLetter(ctx, args[1], args[2], args[3])
		if err == nil && !found {
			return errors.New("dead-letter entry not found")
		}
		return err
	case "replay", "republish":
		name := os.Getenv("SERVICE")
		pool, err := service.OpenDB(ctx, os.Getenv("DATABASE_URL"), name)
		if err != nil {
			return err
		}
		defer pool.Close()
		if args[0] == "republish" {
			return service.Republish(ctx, pool, name)
		}
		if len(args) != 2 {
			return errors.New("usage: portalctl replay <RFC3339>")
		}
		since, err := time.Parse(time.RFC3339, args[1])
		if err != nil {
			return err
		}
		count, err := events.Replay(ctx, pool, since)
		if err == nil {
			fmt.Printf("Queued %d stored events\n", count)
		}
		return err
	default:
		return fmt.Errorf("unknown portalctl command %q", args[0])
	}
}
