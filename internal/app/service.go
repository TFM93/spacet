package app

import (
	"hash/fnv"
	"spacet/internal/app/bookings"
	"spacet/internal/app/healthcheck"
	"spacet/internal/app/spacex"
	"spacet/internal/app/sync"
	"spacet/internal/domain"
	"spacet/pkg/logger"
)

type HealthCheckServiceQueries interface {
	healthcheck.Queries
}
type SpaceXServiceCommands interface {
	spacex.Commands
}
type SpaceXServiceQueries interface {
	spacex.Queries
}

type BookingsServiceCommands interface {
	bookings.Commands
}
type SyncServiceCommands interface {
	sync.Commands
}

// NewHealthCheckQueries creates an instance of HealthCheck Queries that satisfies HealthCheckQueries interface
func NewHealthCheckQueries(repo domain.MonitoringInfraQueries) HealthCheckServiceQueries {
	return healthcheck.NewQueries(repo)
}

// NewSpaceXCommands creates a new instance of SpaceX Commands
func NewSpaceXCommands(logger logger.Interface, client domain.SpaceXAPIQueries, launchPadRepo domain.LaunchPadRepoCommands, launchesRepo domain.LaunchRepoCommands) SpaceXServiceCommands {
	return spacex.NewCommands(logger, client, launchPadRepo, launchesRepo)
}

// NewSpaceXQueries creates a new instance of SpaceX Queries
func NewSpaceXQueries(logger logger.Interface, client domain.SpaceXAPIQueries) SpaceXServiceQueries {
	return spacex.NewQueries(logger, client)
}

// NewBookingsCommands creates a new instance of Booking Commands
func NewBookingsCommands(logger logger.Interface, transaction domain.Transaction, bookingCmds domain.BookingRepoCommands, launchesCmds domain.LaunchRepoCommands, launchesQrs domain.LaunchRepoQueries) BookingsServiceCommands {
	return bookings.NewCommands(logger, transaction, bookingCmds, launchesCmds, launchesQrs)
}

// NewSyncCommands creates a new instance of Sync Commands
func NewSyncCommands(logger logger.Interface, transaction domain.Transaction, repo domain.SyncRepoCommands) SyncServiceCommands {
	return sync.NewCommands(logger, transaction, repo, hashResourceName)
}

func hashResourceName(resourceName string) uint32 {
	// https://pkg.go.dev/hash#pkg-subdirectories
	h := fnv.New32a()
	h.Write([]byte(resourceName))
	return h.Sum32()
}
