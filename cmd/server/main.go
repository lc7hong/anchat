package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/huyng/anchat/internal/server"
)

var (
	dbPath      string
	storageKey  string
	tlsCert    string
	tlsKey     string
	port        string
)

func init() {
	flag.StringVar(&dbPath, "db", "anchat.db", "Path to SQLite database")
	flag.StringVar(&storageKey, "storage-key", "", "Server storage encryption key (required)")
	flag.StringVar(&tlsCert, "cert", "", "TLS certificate path (required)")
	flag.StringVar(&tlsKey, "key", "", "TLS private key path (required)")
	flag.StringVar(&port, "port", ":443", "Port to listen on")
	flag.Parse()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Validate required arguments
	if storageKey == "" {
		log.Fatal("Missing required argument: -storage-key")
	}
	if tlsCert == "" {
		log.Fatal("Missing required argument: -cert")
	}
	if tlsKey == "" {
		log.Fatal("Missing required argument: -key")
	}

	// Check if TLS files exist
	if _, err := os.Stat(tlsCert); os.IsNotExist(err) {
		log.Fatalf("TLS certificate not found: %s", tlsCert)
	}
	if _, err := os.Stat(tlsKey); os.IsNotExist(err) {
		log.Fatalf("TLS private key not found: %s", tlsKey)
	}

	// Create server
	anchat, err := server.New(dbPath, storageKey, tlsCert, tlsKey)
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
	log.Printf("AnChat server starting on %s", port)
	log.Printf("Database: %s", dbPath)
	log.Printf("TLS cert: %s", tlsCert)
	log.Printf("TLS key: %s", tlsKey)

	if err := anchat.Start(port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
