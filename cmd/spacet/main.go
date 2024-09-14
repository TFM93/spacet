package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"spacet/config"
	"syscall"
	"time"

	"spacet/internal/app"
	"spacet/internal/controller/grpc"
	"spacet/internal/controller/http"
	"spacet/internal/infra/postgresql"
	"spacet/internal/infra/spacex"
	"spacet/pkg/grpcserver"
	"spacet/pkg/httpserver"
	"spacet/pkg/logger"
	postgres "spacet/pkg/postgresql"
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

func run(cfg *config.Config, l logger.Interface) error {
	// -------------------------------------------------------------------------
	// Setup Infra

	pg, err := postgres.New(
		cfg.PG.DSN,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.AutoMigrate(true, "../../migrations/postgresql"),
		postgres.WithLogger(l))

	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}
	defer pg.Close()

	txSupplier := postgresql.NewTransactionSupplier(pg)

	lpadCmdRepo := postgresql.NewLaunchPadCommandsRepo(pg, l)
	lauchesCmdRepo := postgresql.NewLaunchesCommandsRepo(pg, l)
	launchesQrRepo := postgresql.NewLaunchesQueriesRepo(pg, l)
	spacexClient := spacex.NewSpaceXQueries(l)

	// -------------------------------------------------------------------------
	// Setup Service Layer

	healthCheckQueries := app.NewHealthCheckQueries(pg)
	spaceXCommands := app.NewSpaceXCommands(l, spacexClient, lpadCmdRepo, lauchesCmdRepo)
	spaceXQueries := app.NewSpaceXQueries(l, spacexClient)
	bookingsCommands := app.NewBookingsCommands(l, txSupplier, postgresql.NewBookingCommandsRepo(pg, l), lauchesCmdRepo, launchesQrRepo)
	syncCommands := app.NewSyncCommands(l, txSupplier, postgresql.NewSyncCommandsRepo(pg, l))

	// updates the launchpad every 30 days
	// TODO: run this as a schedule and configure the launchpad frequency on configs
	if err := syncCommands.SyncIfNecessary(context.Background(), "launchpads", 30*24*time.Hour, spaceXCommands.UpdateLaunchPads); err != nil {
		l.Error("failed to sync launchpads: %s", err)
	}

	bookingsOrchestrator := app.NewBookingsOrchestrator(l, spaceXCommands, spaceXQueries, bookingsCommands, syncCommands)
	fmt.Println(bookingsOrchestrator.SyncOnce(context.Background(), 0))
	// orchestratorInterval := time.Duration(cfg.Orchestrator.Interval) * time.Hour
	// go bookingsOrchestrator.StartScheduledSync(context.Background(), orchestratorInterval)
	// defer bookingsOrchestrator.GracefulStop()

	// -------------------------------------------------------------------------
	// Setup Controller Layer
	httpEngine, err := http.Setup(l, cfg.GRPC.Port, healthCheckQueries)
	if err != nil {
		return fmt.Errorf("httpServer.Setup: %w", err)
	}

	settedUpServer, err := grpc.Setup(l)
	if err != nil {
		return fmt.Errorf("grpcServer.Setup: %w", err)
	}
	// -------------------------------------------------------------------------
	// Start API Servers

	grpcServer := grpcserver.New(settedUpServer, grpcserver.Port(cfg.GRPC.Port), grpcserver.WithLogger(l))
	defer grpcServer.GracefulStop()
	httpServer := httpserver.New(httpEngine, httpserver.Port(cfg.HTTP.Port), httpserver.WithLogger(l))

	// -------------------------------------------------------------------------
	// Shutdown

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("received run signal: " + s.String())
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("run - grpcServer.Notify: %w", err))
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		return fmt.Errorf("httpServer.Shutdown: %w", err)
	}

	return nil
}
