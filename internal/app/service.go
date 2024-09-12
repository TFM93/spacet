package app

import (
	"spacet/internal/app/bookings"
	"spacet/internal/app/healthcheck"
	"spacet/internal/app/spacex"
	"spacet/pkg/logger"
)

// NewHealthCheckQueries creates an instance of HealthCheck Queries that satisfies HealthCheckQueries interface
func NewHealthCheckQueries() healthcheck.Queries {
	return healthcheck.NewQueries()
}

// NewSpaceXCommands creates a new instance of SpaceX Commands
func NewSpaceXCommands(logger logger.Interface) spacex.Commands {
	return spacex.NewCommands(logger)
}

// NewBookingsCommands creates a new instance of Booking Commands
func NewBookingsCommands(logger logger.Interface) bookings.Commands {
	return bookings.NewCommands(logger)
}
