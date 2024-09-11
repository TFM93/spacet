package services

import "context"

// HealthCheckQueries is an interface for checking the health of application dependencies
type HealthCheckQueries interface {
	Check(ctx context.Context) bool
}
