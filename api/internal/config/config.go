// Package config loads service configuration from environment variables.
package config

import (
	"fmt"
	"os"
)

// Config holds all environment-derived settings for the API service.
type Config struct {
	Port                string
	DatabaseURL         string
	AdminBootstrapToken string
}

// Load reads configuration from the environment, applying defaults where
// sensible and failing fast on missing required values.
func Load() (Config, error) {
	cfg := Config{
		Port:                os.Getenv("PORT"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		AdminBootstrapToken: os.Getenv("ADMIN_BOOTSTRAP_TOKEN"),
	}

	if cfg.Port == "" {
		cfg.Port = "8090"
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.AdminBootstrapToken == "" {
		return Config{}, fmt.Errorf("ADMIN_BOOTSTRAP_TOKEN is required (protects POST /api/v1/auth/keys)")
	}

	return cfg, nil
}
