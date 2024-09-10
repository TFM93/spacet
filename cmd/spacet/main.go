package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// -------------------------------------------------------------------------
	// Configuration

	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {

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
	fmt.Println("waiting for interrupt signal")
	<-interrupt
	return nil
}
