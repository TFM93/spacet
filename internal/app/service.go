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

// NewHealthCheckQueries creates an instance of HealthCheck Queries that satisfies HealthCheckQueries interface
func NewHealthCheckQueries() healthcheck.Queries {
	return healthcheck.NewQueries()
}

// NewSpaceXCommands creates a new instance of SpaceX Commands
func NewSpaceXCommands(logger logger.Interface, launchPadRepo domain.LaunchPadRepoCommands, launchesRepo domain.LaunchRepoCommands) spacex.Commands {
	return spacex.NewCommands(logger, launchPadRepo, launchesRepo)
}

// NewBookingsCommands creates a new instance of Booking Commands
func NewBookingsCommands(logger logger.Interface) bookings.Commands {
	return bookings.NewCommands(logger)
}

// NewSyncCommands creates a new instance of Sync Commands
func NewSyncCommands(logger logger.Interface, transaction domain.Transaction, repo domain.SyncRepoCommands) sync.Commands {
	return sync.NewCommands(logger, transaction, repo, hashResourceName)
}

func hashResourceName(resourceName string) uint32 {
	// https://pkg.go.dev/hash#pkg-subdirectories
	h := fnv.New32a()
	h.Write([]byte(resourceName))
	return h.Sum32()
}
