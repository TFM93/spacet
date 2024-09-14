package healthcheck

import (
	"context"
	"spacet/internal/domain"
)

// Queries is an interface for checking the health of application dependencies
type Queries interface {
	Check(ctx context.Context) bool
}

type healthCheckQueries struct {
	repo domain.MonitoringInfraQueries
}

// NewQueries creates a service that satisfies the interface HealthCheckQueries
func NewQueries(repo domain.MonitoringInfraQueries) Queries {
	return &healthCheckQueries{repo: repo}
}

// Check pings the following dependencies:
// - repository
func (h *healthCheckQueries) Check(ctx context.Context) bool {
	return h.repo.Ping(ctx)
}
