package domain

import (
	"context"
)

// MonitoringInfraQueries is an interface for Ping application dependencies
type MonitoringInfraQueries interface {
	// Ping asserts that a given dependency's communication works
	Ping(ctx context.Context) bool
	// IsEnabled checks if the dependency is enabled or not
	IsEnabled() bool
}
