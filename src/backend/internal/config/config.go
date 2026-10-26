package config

import (
	"errors"
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config contains the settings needed to start the API.
type Config struct {
	APIAddr           string
	DatabaseURL       string
	TrustedProxyCIDRs []netip.Prefix
}

// Load reads the local .env without overriding existing environment variables.
func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	cfg := Config{
		APIAddr:     os.Getenv("API_ADDR"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required; copy .env.example to .env or set it in the environment")
	}
	if cfg.APIAddr == "" {
		cfg.APIAddr = "127.0.0.1:8080"
	}
	for _, raw := range strings.Split(os.Getenv("TRUSTED_PROXY_CIDRS"), ",") {
		if strings.TrimSpace(raw) == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(strings.TrimSpace(raw))
		if err != nil {
			return Config{}, fmt.Errorf("TRUSTED_PROXY_CIDRS: %w", err)
		}
		cfg.TrustedProxyCIDRs = append(cfg.TrustedProxyCIDRs, prefix)
	}
	return cfg, nil
}
