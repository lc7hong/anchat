package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds the server configuration
type Config struct {
	Server ServerConfig `toml:"server"`
	DB     DBConfig     `toml:"database"`
	Security SecurityConfig `toml:"security"`
	TLS    TLSConfig    `toml:"tls"`
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port string `toml:"port"`
}

// DBConfig holds database settings
type DBConfig struct {
	Path string `toml:"path"`
}

// SecurityConfig holds security settings
type SecurityConfig struct {
	StorageKey string `toml:"storage_key"`
}

// TLSConfig holds TLS settings
type TLSConfig struct {
	Enabled bool   `toml:"enabled"`
	Cert    string `toml:"cert"`
	Key     string `toml:"key"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: ":8080",
		},
		DB: DBConfig{
			Path: "anchat.db",
		},
		TLS: TLSConfig{
			Enabled: false,
		},
	}
}

// Load loads configuration from a TOML file
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate required fields
	if cfg.Security.StorageKey == "" {
		return nil, fmt.Errorf("missing required field: security.storage_key")
	}

	if cfg.TLS.Enabled {
		if cfg.TLS.Cert == "" {
			return nil, fmt.Errorf("TLS enabled but missing tls.cert")
		}
		if cfg.TLS.Key == "" {
			return nil, fmt.Errorf("TLS enabled but missing tls.key")
		}
	}

	return cfg, nil
}

// Write writes the configuration to a file
func Write(path string, cfg *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
