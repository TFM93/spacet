package spacex

import (
	"context"
	"spacet/internal/domain"
	"spacet/pkg/logger"
)

type Queries interface {
	// GetUpcomingLaunches fetches the upcoming launches from the spacex api
	GetUpcomingLaunches(ctx context.Context) ([]*domain.Launch, error)
}

type queriesHandler struct {
	l            logger.Interface
	spacexClient domain.SpaceXAPIQueries
}

func NewQueries(logger logger.Interface, client domain.SpaceXAPIQueries) Queries {
	return &queriesHandler{l: logger, spacexClient: client}
}

func (h queriesHandler) GetUpcomingLaunches(ctx context.Context) ([]*domain.Launch, error) {
	return h.spacexClient.GetUpcomingLaunches(ctx)
}
