package healthcheck

import (
	"context"
)

// Queries is an interface for checking the health of application dependencies
type Queries interface {
	Check(ctx context.Context) bool
}

type healthCheckQueries struct {
}

// NewQueries creates a service that satisfies the interface HealthCheckQueries
func NewQueries() Queries {
	return &healthCheckQueries{}
}

// Check pings the following dependencies:
// - repository
func (h *healthCheckQueries) Check(ctx context.Context) bool {
	// todo: implement repository healthcheck
	return true
}
