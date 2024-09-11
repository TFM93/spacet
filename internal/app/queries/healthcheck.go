package app

import (
	"context"
	"spacet/internal/app/services"
)

type healthCheckQueries struct {
}

// NewHealthCheckQueries creates a service that satisfies the interface HealthCheckQueries
func NewHealthCheckQueries() services.HealthCheckQueries {
	return &healthCheckQueries{}
}

// Check pings the following dependencies:
// - repository
func (h *healthCheckQueries) Check(ctx context.Context) bool {
	// todo: implement repository healthcheck
	return true
}
