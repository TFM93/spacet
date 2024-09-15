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
	"spacet/internal/app/bookings"
	"spacet/internal/controller/grpc"
	"spacet/internal/controller/http"
	"spacet/internal/domain"
	"spacet/internal/infra/notification"
	"spacet/internal/infra/outbox"
	"spacet/internal/infra/postgresql"
	"spacet/internal/infra/spacex"
	"spacet/pkg/grpcserver"
	"spacet/pkg/httpserver"
	"spacet/pkg/logger"
	postgres "spacet/pkg/postgresql"
	"spacet/pkg/pubsub"
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
	outboxRepoCommands := postgresql.NewOutboxCommandsRepo(pg, l)

	lpadCmdRepo := postgresql.NewLaunchPadCommandsRepo(pg, l)
	lauchesCmdRepo := postgresql.NewLaunchesCommandsRepo(pg, l)
	launchesQrRepo := postgresql.NewLaunchesQueriesRepo(pg, l)
	spacexClient := spacex.NewSpaceXQueries(l)

	pubsubClient, err := pubsub.New(cfg.PubSub.Enabled, cfg.PubSub.ProjectID, pubsub.WithLogger(l))
	if err != nil {
		return fmt.Errorf("pubsub.New: %w", err)
	}

	var notificationService domain.NotificationService
	if cfg.PubSub.Enabled {
		notificationService = notification.NewPubSubNotifierService(pubsubClient, l, cfg.PubSub.LaunchesTopic)
	} else {
		notificationService = notification.NewLoggerNotifierService(l)
	}

	outboxProcessor := outbox.NewProcessor(l, txSupplier, outboxRepoCommands, notificationService)
	interval := time.Duration(cfg.Notifications.Interval) * time.Second
	go outboxProcessor.StartScheduleProcess(context.Background(), interval, cfg.Notifications.MaxBatchSize)

	// -------------------------------------------------------------------------
	// Setup Service Layer

	healthCheckQueries := app.NewHealthCheckQueries(pg)
	spaceXCommands := app.NewSpaceXCommands(l, spacexClient, lpadCmdRepo, lauchesCmdRepo)
	spaceXQueries := app.NewSpaceXQueries(l, spacexClient)
	bookingsCommands := app.NewBookingsCommands(l, txSupplier, postgresql.NewBookingCommandsRepo(pg, l), lauchesCmdRepo, launchesQrRepo)
	bookingsQueries := app.NewBookingsQueries(l, txSupplier, postgresql.NewBookingQueriesRepo(pg, l))
	syncCommands := app.NewSyncCommands(l, txSupplier, postgresql.NewSyncCommandsRepo(pg, l))

	// updates the launchpad every 30 days
	// TODO: run this as a schedule and configure the launchpad frequency on configs
	if err := syncCommands.SyncIfNecessary(context.Background(), "launchpads", 30*24*time.Hour, spaceXCommands.UpdateLaunchPads); err != nil {
		l.Error("failed to sync launchpads: %s", err)
	}

	fmt.Println(bookingsCommands.BookALaunch(context.Background(), bookings.BookALaunchReq{
		LaunchpadID: "5e9e4501f509094ba4566f84",
		Date:        time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
		Destination: domain.DestinationAsteroidBelt,
		FirstName:   "sd",
		LastName:    "ds",
		Gender:      domain.GenderFemale,
		Birthday:    time.Date(1990, 12, 01, 0, 0, 0, 0, time.UTC),
	}))
	bookingsOrchestrator := app.NewBookingsOrchestrator(l, spaceXCommands, spaceXQueries, bookingsCommands, syncCommands, outboxRepoCommands)
	orchestratorInterval := time.Duration(cfg.Orchestrator.Interval) * time.Hour
	fmt.Sprintln(bookingsOrchestrator.SyncOnce(context.Background(), 0))
	go bookingsOrchestrator.StartScheduledSync(context.Background(), orchestratorInterval)
	defer bookingsOrchestrator.GracefulStop()

	// -------------------------------------------------------------------------
	// Setup Controller Layer
	httpEngine, err := http.Setup(l, cfg.GRPC.Port, healthCheckQueries)
	if err != nil {
		return fmt.Errorf("httpServer.Setup: %w", err)
	}

	settedUpServer, err := grpc.Setup(l, bookingsCommands, bookingsQueries)
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
