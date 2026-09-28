// Package config loads the explicitly supported Serega process configuration.
package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	defaultHTTPAddr        = ":8080"
	defaultShutdownTimeout = 15 * time.Second
)

// Load reads environment variables. A serving process always needs a database;
// migrations intentionally use their own, isolated configuration path.
func Load() (Config, error) {
	return FromLookup(os.LookupEnv)
}

// FromLookup makes configuration validation testable without mutating process
// environment variables.
func FromLookup(lookup func(string) (string, bool)) (Config, error) {
	value := func(key, fallback string) string {
		if v, ok := lookup(key); ok {
			return strings.TrimSpace(v)
		}

		return fallback
	}

	cfg := Config{
		HTTPAddr:        value("SEREGA_HTTP_ADDR", defaultHTTPAddr),
		DatabaseURL:     value("DATABASE_URL", ""),
		LogFormat:       value("SEREGA_LOG_FORMAT", "json"),
		ShutdownTimeout: defaultShutdownTimeout,
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required when serving Serega")
	}

	if _, _, err := net.SplitHostPort(cfg.HTTPAddr); err != nil {
		return Config{}, fmt.Errorf("SEREGA_HTTP_ADDR must be host:port: %w", err)
	}

	if cfg.LogFormat != "json" && cfg.LogFormat != "text" {
		return Config{}, errors.New("SEREGA_LOG_FORMAT must be json or text")
	}

	if raw := value("SEREGA_SHUTDOWN_TIMEOUT", ""); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil || d <= 0 {
			return Config{}, errors.New("SEREGA_SHUTDOWN_TIMEOUT must be a positive duration")
		}

		cfg.ShutdownTimeout = d
	}

	proxies, err := parseCIDRs(value("SEREGA_TRUSTED_PROXY_CIDRS", ""))
	if err != nil {
		return Config{}, err
	}

	cfg.TrustedProxies = proxies

	secure := value("SEREGA_COOKIE_SECURE", "true")
	if secure != "true" && secure != "false" {
		return Config{}, errors.New("SEREGA_COOKIE_SECURE must be true or false")
	}

	cfg.SecureCookies = secure == "true"

	return cfg, nil
}
