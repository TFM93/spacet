package spacex

import (
	"context"
	"spacet/internal/domain"
	"spacet/pkg/logger"
)

type Queries interface {
	GetUpcomingLaunches(ctx context.Context) ([]*domain.Launch, error)
}

type queriesHandler struct {
	l            logger.Interface
	spacexClient domain.SpaceXAPIQueries
}

func NewQueries(logger logger.Interface, client domain.SpaceXAPIQueries) Queries {
	return &queriesHandler{l: logger, spacexClient: client}
}

// GetUpcomingLaunches ... todo describe.
func (h queriesHandler) GetUpcomingLaunches(ctx context.Context) ([]*domain.Launch, error) {
	return h.spacexClient.GetUpcomingLaunches(ctx)
}
