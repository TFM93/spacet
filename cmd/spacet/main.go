package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"spacet/config"
	"syscall"

	"spacet/pkg/logger"
)

func main() {
	// -------------------------------------------------------------------------
	// Configuration

	configPath := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	l := logger.New(cfg.LogLevel)

	if err := run(cfg, l); err != nil {
		l.Error("Run error: %s", err)
		os.Exit(1)
	}
}

func run(_ *config.Config, l logger.Interface) error {

	// -------------------------------------------------------------------------
	// Setup Infra

	// -------------------------------------------------------------------------
	// Setup Service Layer

	// -------------------------------------------------------------------------
	// Setup Controller Layer

	// -------------------------------------------------------------------------
	// Start API Servers

	// -------------------------------------------------------------------------
	// Shutdown

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	l.Info("waiting for interrupt signal")
	<-interrupt
	return nil
}
