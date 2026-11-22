package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Migrate stops concurrent migration jobs with one session lock spanning the
// legacy schema, demo data, split schema and role setup.
func Migrate(ctx context.Context, adminURL, dir string, seed bool, roleURLs map[string]string) error {
	cfg, err := pgx.ParseConfig(adminURL)
	if err != nil || adminURL == "" {
		return fmt.Errorf("invalid PORTAL67_ADMIN_DATABASE_URL")
	}
	cfg.RuntimeParams["search_path"] = "public"
	delete(cfg.RuntimeParams, "options")
	db := stdlib.OpenDB(*cfg)
	defer db.Close()
	lock, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer lock.Close()
	if _, err := lock.ExecContext(ctx, "SELECT pg_advisory_lock(67220001)"); err != nil {
		return err
	}

	up := func(part, table string) error {
		provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS(filepath.Join(dir, part)), goose.WithTableName(table))
		if err != nil {
			return err
		}
		_, err = provider.Up(ctx)
		return err
	}
	if err := up("migrations", "goose_db_version"); err != nil {
		return fmt.Errorf("migrate legacy schema: %w", err)
	}
	if seed {
		var initialize bool
		err := db.QueryRowContext(ctx, `
			SELECT to_regnamespace('identity') IS NULL AND
				(to_regclass('public.portal67_demo_migrations') IS NOT NULL OR NOT EXISTS (SELECT 1 FROM public.users))
		`).Scan(&initialize)
		if err != nil {
			return err
		}
		if initialize {
			if err := up("seeds", "portal67_demo_migrations"); err != nil {
				return fmt.Errorf("seed demo data: %w", err)
			}
		}
	}
	if err := up("services", "portal67_services_migrations"); err != nil {
		return fmt.Errorf("migrate service schemas: %w", err)
	}

	roles, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer roles.Rollback()
	for service, rawURL := range roleURLs {
		role := "portal67_" + service
		credentials, err := pgx.ParseConfig(rawURL)
		if err != nil || credentials.User != role || credentials.Password == "" || rawURL == "" {
			return fmt.Errorf("%s_DATABASE_URL must include credentials for %s", strings.ToUpper(service), role)
		}
		var query string
		if err := roles.QueryRowContext(ctx, "SELECT format('ALTER ROLE %I LOGIN PASSWORD %L', $1::text, $2::text)",
			role, credentials.Password).Scan(&query); err != nil {
			return fmt.Errorf("configure role %s: %w", role, err)
		}
		if _, err := roles.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("configure role %s: %w", role, err)
		}
	}
	return roles.Commit()
}
