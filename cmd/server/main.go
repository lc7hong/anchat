package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/huyng/anchat/internal/config"
	"github.com/huyng/anchat/internal/server"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "anchat.toml", "Path to config file")
	flag.Parse()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load configuration (optional)
	var cfg *config.Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, use defaults
		cfg = config.DefaultConfig()
		log.Printf("No config file found at %s, using defaults", configPath)
		log.Println("Create a config file to customize settings (see anchat.example.toml)")
	} else {
		var err error
		cfg, err = config.Load(configPath)
		if err != nil {
			log.Fatalf("Failed to load config from %s: %v", configPath, err)
		}
		log.Printf("Loaded config from %s", configPath)
	}

	// Create server
	anchat, err := server.New(cfg.DB.Path, cfg.Security.StorageKey, cfg.TLS.Cert, cfg.TLS.Key, cfg.TLS.Enabled)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer anchat.Stop(context.Background())

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		anchat.Stop(context.Background())
		os.Exit(0)
	}()

	// Start server
	log.Printf("AnChat server starting on %s", cfg.Server.Port)
	log.Printf("Database: %s", cfg.DB.Path)
	if cfg.TLS.Enabled {
		log.Printf("TLS enabled: cert=%s, key=%s", cfg.TLS.Cert, cfg.TLS.Key)
	} else {
		log.Println("TLS disabled - running in plaintext mode")
		log.Println("WARNING: Session tokens and public keys will be sent unencrypted")
		log.Println("Messages remain end-to-end encrypted")
	}

	if err := anchat.Start(cfg.Server.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
