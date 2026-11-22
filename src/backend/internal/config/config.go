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
	RedisURL          string
	RedisKeyPrefix    string
	TrustedProxyCIDRs []netip.Prefix
}

// Load reads the local .env without overriding existing environment variables.
func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	cfg := Config{
		APIAddr:        os.Getenv("API_ADDR"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		RedisURL:       os.Getenv("REDIS_URL"),
		RedisKeyPrefix: strings.TrimSpace(os.Getenv("REDIS_KEY_PREFIX")),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required; copy .env.example to .env or set it in the environment")
	}
	if cfg.APIAddr == "" {
		cfg.APIAddr = "127.0.0.1:8080"
	}
	trusted, err := TrustedProxyCIDRs()
	cfg.TrustedProxyCIDRs = trusted
	return cfg, err
}

func TrustedProxyCIDRs() ([]netip.Prefix, error) {
	var cidrs []netip.Prefix
	for _, raw := range strings.Split(os.Getenv("TRUSTED_PROXY_CIDRS"), ",") {
		if strings.TrimSpace(raw) == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(strings.TrimSpace(raw))
		if err != nil {
			return nil, fmt.Errorf("TRUSTED_PROXY_CIDRS: %w", err)
		}
		cidrs = append(cidrs, prefix)
	}
	return cidrs, nil
}
